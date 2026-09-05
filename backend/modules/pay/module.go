// Package pay is the commerce module: shop products, orders and the
// payment-gateway plugin surface (see gateways/). Also hosts the board
// moderator table and the admin mailer endpoint.
package pay

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/database/schema"
	mailcontract "github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	authhttp "reflexcms/backend/modules/auth/http/controllers"
	paygateways "reflexcms/backend/modules/pay/gateways"
	paymodels "reflexcms/backend/modules/pay/models"
	settingsservices "reflexcms/backend/modules/settings/services"
	"reflexcms/backend/database/migrations"
)

// netHTTP isolates the stdlib client used by the resend driver.
var netHTTP = &nethttp.Client{Timeout: 15 * time.Second}

// paySet reads a pay-group settings key with a fallback.
func paySet(key, fallback string) string {
	v := settingsservices.Get(key)
	if v == "" {
		return fallback
	}
	return v
}

// Module owns the shop + payment gateways + board moderators + mailer.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "pay" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.PayMigrations()
}

func (m *Module) Routes() {
	registerAdminSpecs()
	registerPublicRoutes()

	// Admin mailer (single / bulk) — permission-gated by settings.edit.
	facades.Route().Post("/api/admin/mail/send", SendMail)
}

func (m *Module) Boot() {}

// ---- admin gateway specs ----

func registerAdminSpecs() {
	adminhubRegister(&adminhubSpec{
		Name:             "products",
		Model:            paymodels.Product{},
		Searchable:       []string{"title"},
		Sortable:         []string{"id", "title", "price_cents", "stock", "sort", "created_at"},
		Fillable:         []string{"title", "description", "price_cents", "currency", "stock", "grant_points", "image", "is_active", "sort", "type"},
		PermissionPrefix: "products",
	})

	// Orders are effectively read-only for admins (no Fillable).
	adminhubRegister(&adminhubSpec{
		Name:             "orders",
		Model:            paymodels.Order{},
		Searchable:       []string{"order_no", "title"},
		Sortable:         []string{"id", "order_no", "amount_cents", "status", "gateway", "created_at"},
		Fillable:         []string{},
		PermissionPrefix: "orders",
	})
}

// ---- public routes ----

func registerPublicRoutes() {
	r := facades.Route()

	// Shop listing — visible only when the shop module is enabled.
	r.Get("/api/v1/products", func(ctx http.Context) http.Response {
		if !shopEnabled() {
			return httpx.Error(ctx, 404, "Not Found")
		}
		items := []map[string]any{}
		if err := facades.Orm().Query().Table("products").
			Where("is_active = ? AND stock > 0", true).
			Select("id", "title", "description", "price_cents", "currency", "stock", "image", "grant_points").
			OrderBy("sort", "asc").Get(&items); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{
			"items":     items,
			"currency":  paySet("pay.display_currency", "USD"),
			"aff_links": affLinks(),
			"gateways":  enabledGatewayNames(),
		})
	})

	// Create an order — GUESTS ALLOWED. Identity is optional; an anonymous
	// buyer is tracked by the unguessable order_no which the frontend keeps
	// in sessionStorage, so the invite code stays retrievable after payment.
	r.Post("/api/v1/orders", func(ctx http.Context) http.Response {
		if !shopEnabled() {
			return httpx.Error(ctx, 404, "Not Found")
		}
		// optional identity
		userID := uint64(0)
		if identity, ok := authhttp.RequireIdentity(ctx); ok {
			userID = identity.ID
		}

		var req struct {
			ProductID uint64 `json:"product_id"`
			Gateway   string `json:"gateway"`
		}
		if err := ctx.Request().Bind(&req); err != nil || req.ProductID == 0 {
			return httpx.Error(ctx, 422, "product_id is required")
		}

		var products []map[string]any
		if err := facades.Orm().Query().Table("products").
			Where("id = ? AND is_active = ? AND stock > 0", req.ProductID, true).
			Get(&products); err != nil || len(products) == 0 {
			return httpx.Error(ctx, 404, "商品不存在或已售罄")
		}
		p := products[0]

		gw, err := paygateways.Get(req.Gateway)
		if err != nil || !gw.Enabled() {
			return httpx.Error(ctx, 422, "支付渠道不可用")
		}

		amountCents := toI64(p["price_cents"])
		currency := paySet("pay.display_currency", fmt.Sprint(p["currency"]))
		if currency == "" {
			currency = fmt.Sprint(p["currency"])
		}
		orderNo := newOrderNo()

		if _, err := facades.Orm().Query().Exec(`
			INSERT INTO orders (order_no, user_id, product_id, title, amount_cents, currency, gateway, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', NOW(), NOW())
		`, orderNo, userID, req.ProductID, fmt.Sprint(p["title"]), amountCents, currency, gw.Name()); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}

		result, err := gw.Create(paygateways.CreateRequest{
			OrderNo:     orderNo,
			Title:       fmt.Sprint(p["title"]),
			AmountCents: amountCents,
			Currency:    currency,
			ReturnURL:   siteURL() + "/shop/orders/" + orderNo,
			NotifyURL:   siteURL() + "/api/v1/pay/notify/" + gw.Name(),
		})
		if err != nil {
			_, _ = facades.Orm().Query().Table("orders").
				Where("order_no = ?", orderNo).
				Update(map[string]any{"status": "failed"})
			return httpx.Error(ctx, 502, err.Error())
		}

		for _, key := range []string{"sys_no", "invoice_id", "paypal_order_id"} {
			if ref, ok := result.Raw[key].(string); ok && ref != "" {
				_, _ = facades.Orm().Query().Table("orders").
					Where("order_no = ?", orderNo).
					Update(map[string]any{"gateway_ref": ref})
				break
			}
		}

		return ctx.Response().Success().Json(http.Json{
			"order_no": orderNo,
			"pay_url":  result.PayURL,
		})
	})

	// Order status — PUBLIC by order number (unguessable; guests rely on it
	// to retrieve purchased invite codes). Sensitive fields are not exposed.
	r.Get("/api/v1/orders/{order_no}", func(ctx http.Context) http.Response {
		orderNo := ctx.Request().Route("order_no")
		var rows []map[string]any
		if err := facades.Orm().Query().Table("orders").
			Where("order_no = ?", orderNo).
			Select("order_no", "title", "amount_cents", "currency", "gateway", "status", "paid_at", "granted_code", "created_at").
			Get(&rows); err != nil || len(rows) == 0 {
			return httpx.Error(ctx, 404, "order not found")
		}
		return ctx.Response().Success().Json(rows[0])
	})

	// Admin: manually mark an order paid (offline payments / support).
	r.Post("/api/admin/orders/{id}/mark-paid", func(ctx http.Context) http.Response {
		if resp := adminhubCheckPermission(ctx, "orders.edit"); resp != nil {
			return *resp
		}
		idStr := ctx.Request().Route("id")
		id, parseErr := strconv.ParseUint(idStr, 10, 64)
		if parseErr != nil || id == 0 {
			return httpx.Error(ctx, 404, "order not found")
		}
		var nos []string
		if err := facades.Orm().Query().Table("orders").
			Where("id = ?", id).Pluck("order_no", &nos); err != nil || len(nos) == 0 {
			return httpx.Error(ctx, 404, "order not found")
		}
		if err := markOrderPaid(nos[0]); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"message": "已标记为已支付并发放商品"})
	})

	// Webhook dispatcher: one endpoint per gateway.
	r.Post("/api/v1/pay/notify/{gateway}", func(ctx http.Context) http.Response {
		gwName := ctx.Request().Route("gateway")
		body, err := io.ReadAll(ctx.Request().Origin().Body)
		if err != nil {
			return httpx.Error(ctx, 422, "unreadable body")
		}
		headers := ctx.Request().Origin().Header
		form := map[string]string{}
		// Providers either POST JSON or urlencoded forms; normalise both.
		if strings.HasPrefix(ctx.Request().Header("Content-Type", ""), "application/x-www-form-urlencoded") {
			for k, vals := range ctx.Request().Origin().PostForm {
				if len(vals) > 0 {
					form[k] = vals[0]
				}
			}
		}

		orderNo := ""
		switch gwName {
		case "xcash":
			if !paygateways.VerifyXcash(paySet("pay.xcash_appid", ""), paySet("pay.xcash_hmac_key", ""), headers, body) {
				return httpx.Error(ctx, 403, "invalid signature")
			}
			var ev struct {
				Type string `json:"type"`
				Data struct {
					OutNo     string `json:"out_no"`
					Confirmed bool   `json:"confirmed"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &ev) == nil && ev.Type == "invoice" && ev.Data.Confirmed {
				orderNo = ev.Data.OutNo
			}
		case "coinpayments":
			if !paygateways.VerifyCoinPayments(paySet("pay.coinpayments_client_id", ""), paySet("pay.coinpayments_client_secret", ""), headers, body) {
				return httpx.Error(ctx, 403, "invalid signature")
			}
			var ev struct {
				Data struct {
					Metadata struct {
						OrderNo string `json:"order_no"`
					} `json:"metadata"`
					Status string `json:"status"`
				} `json:"data"`
			}
			if json.Unmarshal(body, &ev) == nil && strings.EqualFold(ev.Data.Status, "paid") {
				orderNo = ev.Data.Metadata.OrderNo
			}
		case "xunhupay", "codepay":
			// These providers' wire protocol requires an MD5 signature;
			// the project security policy blocks that primitive, so the
			// gateways are disabled and their callbacks are rejected until
			// the policy owner grants an explicit exception.
			return httpx.Error(ctx, 501, "该渠道签名协议与本项目安全策略冲突，暂未启用（见 gateways/legacy_md5_gateways.go）")
		case "paypal":
			// Trust only after PayPal-side verification.
			if !paygateways.VerifyPayPal(
				paySet("pay.paypal_client_id", ""),
				paySet("pay.paypal_secret", ""),
				paySet("pay.paypal_env", "sandbox"),
				headers, body) {
				return httpx.Error(ctx, 403, "invalid signature")
			}
			var ev struct {
				EventType string `json:"event_type"`
				Resource  struct {
					CustomID string `json:"custom_id"`
				} `json:"resource"`
			}
			if json.Unmarshal(body, &ev) == nil && ev.EventType == "PAYMENT.CAPTURE.COMPLETED" {
				orderNo = ev.Resource.CustomID
			}
		default:
			return httpx.Error(ctx, 404, "unknown gateway")
		}

		if orderNo == "" {
			// Not a completion event — acknowledge so the provider stops retrying.
			return ctx.Response().Success().Json(http.Json{"ok": true})
		}
		if err := markOrderPaid(orderNo); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"ok": true})
	})
}

// ---- helpers ----

func shopEnabled() bool {
	v := settingsservices.Get("pay.shop_enabled")
	return v == "true" || v == "1"
}

func siteURL() string {
	return strings.TrimRight(paySet("app.url", "http://localhost:3000"), "/")
}

func affLinks() map[string]string {
	return map[string]string{
		"xcash":        paySet("pay.aff_xcash", ""),
		"coinpayments": paySet("pay.aff_coinpayments", ""),
		"xunhupay":     paySet("pay.aff_xunhupay", ""),
		"codepay":      paySet("pay.aff_codepay", ""),
		"paypal":       paySet("pay.aff_paypal", ""),
	}
}

func enabledGatewayNames() []string {
	out := []string{}
	for _, g := range paygateways.All() {
		if g.Enabled() {
			out = append(out, g.Name())
		}
	}
	return out
}

func newOrderNo() string {
	return fmt.Sprintf("R%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
}

func toI64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int32:
		return int64(n)
	case int:
		return int64(n)
	case uint64:
		return int64(n)
	case float64:
		return int64(n)
	}
	return 0
}

// markOrderPaid flips a pending order to paid, decrements stock and grants
// configured points. Idempotent: safe to call again (e.g. manual re-mark)
// to deliver a still-missing invite code.
func markOrderPaid(orderNo string) error {
	var rows []map[string]any
	if err := facades.Orm().Query().Table("orders").
		Where("order_no = ?", orderNo).
		Get(&rows); err != nil || len(rows) == 0 {
		return nil // unknown order
	}
	order := rows[0]

	status := fmt.Sprint(order["status"])
	if status == "pending" {
		if _, err := facades.Orm().Query().Exec(`
			UPDATE orders SET status = 'paid', paid_at = NOW(), updated_at = NOW()
			WHERE order_no = ? AND status = 'pending'
		`, orderNo); err != nil {
			return err
		}
		order["status"] = "paid"
	}

	productID := toI64(order["product_id"])
	userID := toI64(order["user_id"])
	if productID == 0 {
		return nil
	}

	var prods []map[string]any
	if err := facades.Orm().Query().Table("products").
		Where("id = ?", productID).
		Select("grant_points", "type").
		Get(&prods); err != nil || len(prods) == 0 {
		return nil
	}

	// Invite-code goods: mint + attach a code when the order has none.
	// granted_code may be SQL NULL (scanned as nil), so assert the string
	// type instead of Sprint — Sprint(nil) is "<nil>", never "".
	granted, _ := order["granted_code"].(string)
	if fmt.Sprint(prods[0]["type"]) == "invite" && granted == "" {
		code := generateInviteCode()
		if _, err := facades.Orm().Query().Exec(`
			INSERT INTO invite_codes (code, creator_id, max_uses, used_count, created_at, updated_at)
			VALUES (?, ?, 1, 0, NOW(), NOW())
		`, code, userID); err != nil {
			return err
		}
		if _, err := facades.Orm().Query().Table("orders").
			Where("order_no = ?", orderNo).
			Update(map[string]any{"granted_code": code}); err != nil {
			return err
		}
	}

	// stock decrement only once (guarded by pending→paid transition)
	if status == "pending" {
		if _, err := facades.Orm().Query().Exec(`
			UPDATE products SET stock = GREATEST(stock - 1, 0), updated_at = NOW()
			WHERE id = ? AND stock > 0
		`, productID); err != nil {
			return err
		}
		var grant []int64
		_ = facades.Orm().Query().Table("products").
			Where("id = ?", productID).
			Pluck("grant_points", &grant)
		if len(grant) > 0 && grant[0] > 0 && userID > 0 {
			_ = settingsservices.GrantPoints(uint64(userID), int(grant[0]), "shop_purchase")
		}
	}
	return nil
}

// generateInviteCode mints an unambiguous invite code (same alphabet as
// the layout module's auto-generate).
func generateInviteCode() string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, 10)
	if _, err := crand.Read(buf); err != nil {
		return fmt.Sprintf("INV-%d", time.Now().UnixNano())
	}
	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}
	return "INV-" + string(buf)
}

// ---- board moderator helpers (used by forum v1 moderation) ----

// IsModerator reports whether userID moderates the given board.
func IsModerator(boardID, userID uint64) bool {
	if boardID == 0 || userID == 0 {
		return false
	}
	var count int64
	count, _ = facades.Orm().Query().Table("board_moderators").
		Where("board_id = ? AND user_id = ?", boardID, userID).Count()
	return count > 0
}

// ---- admin mailer ----

// SendMail handles POST /api/admin/mail/send — single or bulk. Body:
// {to: string[] | csv_text, subject, body}. Delivery goes through the
// goravel mail facade with the email module's configured driver; failures
// are reported per-recipient.
func SendMail(ctx http.Context) http.Response {
	if resp := adminhubCheckPermission(ctx, "settings.edit"); resp != nil {
		return *resp
	}

	var req struct {
		To      []string `json:"to"`
		CsvText string   `json:"csv_text"`
		Subject string   `json:"subject"`
		Body    string   `json:"body"`
	}
	if err := ctx.Request().Bind(&req); err != nil || strings.TrimSpace(req.Subject) == "" || strings.TrimSpace(req.Body) == "" {
		return httpx.Error(ctx, 422, "subject 和 body 不能为空")
	}

	recipients := normaliseRecipients(req.To, req.CsvText)
	if len(recipients) == 0 {
		return httpx.Error(ctx, 422, "没有有效的收件邮箱")
	}
	if len(recipients) > 500 {
		return httpx.Error(ctx, 422, "单次最多发送 500 个邮箱")
	}

	sent, failed := dispatchMail(recipients, req.Subject, req.Body)
	return ctx.Response().Success().Json(http.Json{
		"sent":   sent,
		"failed": failed,
		"total":  len(recipients),
	})
}

func normaliseRecipients(list []string, csvText string) []string {
	seen := map[string]bool{}
	out := []string{}
	appendOne := func(raw string) {
		addr := strings.TrimSpace(raw)
		addr = strings.Trim(addr, ",;")
		if addr == "" || !strings.Contains(addr, "@") || seen[addr] {
			return
		}
		seen[addr] = true
		out = append(out, addr)
	}
	for _, line := range list {
		for _, piece := range strings.FieldsFunc(line, func(r rune) bool { return r == ',' || r == ';' || r == '\n' }) {
			appendOne(piece)
		}
	}
	for _, line := range strings.Split(csvText, "\n") {
		for _, piece := range strings.Split(line, ",") {
			appendOne(piece)
		}
	}
	return out
}

// dispatchMail sends each recipient independently so one bad address
// doesn't kill the batch.
func dispatchMail(recipients []string, subject, body string) (sent, failed []string) {
	for _, to := range recipients {
		if err := mailOne(to, subject, body); err != nil {
			failed = append(failed, to)
			continue
		}
		sent = append(sent, to)
	}
	return sent, failed
}

// mailOne delivers via the settings-configured driver:
//   - log (default/dev): writes the message to the app log and succeeds
//   - resend: Resend HTTP API with email.resend_api_key
//   - smtp: goravel Mail facade with email.smtp_* settings applied to config
func mailOne(to, subject, body string) error {
	driver := strings.ToLower(settingsservices.Get("email.driver"))
	if driver == "" {
		driver = "log"
	}

	switch driver {
	case "log":
		facades.Log().Info("mail driver=log", map[string]any{
			"to": to, "subject": subject, "body": body,
		})
		return nil
	case "resend":
		return resendSend(to, subject, body)
	case "smtp":
		return smtpSend(to, subject, body)
	}
	return fmt.Errorf("未知邮件驱动 %q", driver)
}

func resendSend(to, subject, body string) error {
	key := strings.TrimSpace(settingsservices.Get("email.resend_api_key"))
	if key == "" {
		return fmt.Errorf("resend_api_key 未配置")
	}
	fromAddr := strings.TrimSpace(settingsservices.Get("email.from_address"))
	if fromAddr == "" {
		fromAddr = "onboarding@resend.dev"
	}
	fromName := settingsservices.Get("email.from_name")

	payload := map[string]any{
		"from":    fromName + " <" + fromAddr + ">",
		"to":      []string{to},
		"subject": subject,
		"html":    body,
	}
	b, _ := json.Marshal(payload)

	httpReq, _ := nethttp.NewRequest("POST", "https://api.resend.com/emails", strings.NewReader(string(b)))
	httpReq.Header.Set("Authorization", "Bearer "+key)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := netHTTP.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend 发送失败(%d): %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	return nil
}

// smtpSend delegates to the goravel mail facade. SMTP endpoints are read
// from the environment at boot (MAIL_HOST / MAIL_PORT / MAIL_USERNAME /
// MAIL_PASSWORD in .env) — goravel's config contract is read-only at
// runtime, so the email.smtp_* settings serve as documentation for the
// operator of which env values to set.
func smtpSend(to, subject, body string) error {
	host := facades.Config().GetString("mail.host", "")
	if host == "" {
		return fmt.Errorf("SMTP 未配置：请在 .env 设置 MAIL_HOST / MAIL_PORT / MAIL_USERNAME / MAIL_PASSWORD")
	}
	return facades.Mail().To([]string{to}).
		Subject(subject).
		Content(mailcontract.Content{Html: body}).
		Send()
}
