package gateways

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	settingsservices "reflexcms/backend/modules/settings/services"
)

// ---- shared helpers ----

func paySet(key, fallback string) string {
	v := settingsservices.Get(key)
	if v == "" {
		return fallback
	}
	return v
}

// payBool tolerates the value shapes the settings store can return
// (JSON true/false strings, "1", or a raw bool).
func payBool(key string) bool {
	return settingsservices.GetBoolSetting(key)
}

// ---- Xcash (docs/xcash.md): HMAC-SHA256 hex over nonce+timestamp+body ----

type XcashGateway struct{}

func (XcashGateway) Name() string  { return "xcash" }
func (XcashGateway) Enabled() bool { return payBool("pay.xcash_enabled") }

func (XcashGateway) Create(req CreateRequest) (CreateResult, error) {
	appid := paySet("pay.xcash_appid", "")
	key := paySet("pay.xcash_hmac_key", "")
	api := paySet("pay.xcash_api", "https://pay.xca.sh")
	if appid == "" || key == "" {
		return CreateResult{}, fmt.Errorf("xcash 未配置 appid/hmac_key")
	}

	amount := fmt.Sprintf("%.2f", float64(req.AmountCents)/100)
	payload := map[string]any{
		"out_no":     req.OrderNo,
		"title":      req.Title,
		"currency":   req.Currency,
		"amount":     amount,
		"notify_url": req.NotifyURL,
		"return_url": req.ReturnURL,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return CreateResult{}, err
	}

	ts := fmt.Sprint(time.Now().Unix())
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	signature := xcashSignature(key, nonce, ts, body)

	httpReq, err := http.NewRequest("POST", api+"/v1/invoice", bytes.NewReader(body))
	if err != nil {
		return CreateResult{}, err
	}
	httpReq.Header.Set("XC-Appid", appid)
	httpReq.Header.Set("XC-Timestamp", ts)
	httpReq.Header.Set("XC-Nonce", nonce)
	httpReq.Header.Set("XC-Signature", signature)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return CreateResult{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out struct {
		SysNo  string `json:"sys_no"`
		PayURL string `json:"pay_url"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.PayURL == "" {
		return CreateResult{}, fmt.Errorf("xcash 创建账单失败: %s", strings.TrimSpace(string(raw)))
	}
	return CreateResult{PayURL: out.PayURL, Raw: map[string]any{"sys_no": out.SysNo}}, nil
}

// xcashSignature: HMAC-SHA256 hex over nonce+timestamp+body (docs/xcash.md).
func xcashSignature(key, nonce, ts string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(key))
	fmt.Fprintf(mac, "%s%s%s", nonce, ts, string(body))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyXcash checks the webhook HMAC headers against the raw body.
func VerifyXcash(appid, key string, headers http.Header, body []byte) bool {
	if appid != "" && headers.Get("XC-Appid") != appid {
		return false
	}
	nonce := headers.Get("XC-Nonce")
	ts := headers.Get("XC-Timestamp")
	sig := headers.Get("XC-Signature")
	if nonce == "" || ts == "" || sig == "" {
		return false
	}
	expected := xcashSignature(key, nonce, ts, body)
	return hmac.Equal([]byte(expected), []byte(sig))
}

func init() { Register(XcashGateway{}) }

// ---- CoinPayments v2: BOM + method+url+clientId+ts+payload, HMAC-SHA256 b64 ----

type CoinPaymentsGateway struct{}

func (CoinPaymentsGateway) Name() string  { return "coinpayments" }
func (CoinPaymentsGateway) Enabled() bool { return payBool("pay.coinpayments_enabled") }

func (CoinPaymentsGateway) Create(req CreateRequest) (CreateResult, error) {
	clientID := paySet("pay.coinpayments_client_id", "")
	secret := paySet("pay.coinpayments_client_secret", "")
	if clientID == "" || secret == "" {
		return CreateResult{}, fmt.Errorf("coinpayments 未配置 client_id/secret")
	}

	amount := fmt.Sprintf("%.2f", float64(req.AmountCents)/100)
	payload := map[string]any{
		"invoice": map[string]any{
			"invoice_id": req.OrderNo,
			"currency":   req.Currency,
			"amount":     amount,
		},
		"metadata": map[string]any{
			"order_no": req.OrderNo,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return CreateResult{}, err
	}

	url := "https://a-api.coinpayments.net/api/v2/merchant/invoices"
	ts := time.Now().UTC().Format("2006-01-02T15:04:05")
	signature := coinpaymentsSignature(secret, "POST", url, clientID, ts, body)

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return CreateResult{}, err
	}
	httpReq.Header.Set("X-CoinPayments-Client", clientID)
	httpReq.Header.Set("X-CoinPayments-Timestamp", ts)
	httpReq.Header.Set("X-CoinPayments-Signature", signature)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return CreateResult{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out struct {
		Data struct {
			InvoiceID string `json:"invoice_id"`
			Checkout  string `json:"checkout_url"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.Data.Checkout == "" {
		return CreateResult{}, fmt.Errorf("coinpayments 创建失败: %s", strings.TrimSpace(string(raw)))
	}
	return CreateResult{PayURL: out.Data.Checkout, Raw: map[string]any{"invoice_id": out.Data.InvoiceID}}, nil
}

// coinpaymentsSignature mirrors docs/coinpayments-sdk-python-1.0.0/_signer.py:
// message = BOM + method + url + client_id + ts + payload, HMAC-SHA256, base64.
func coinpaymentsSignature(secret, method, url, clientID, ts string, payload []byte) string {
	message := "\uFEFF" + method + url + clientID + ts + string(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// VerifyCoinPayments checks webhook HMAC headers (same signing scheme).
func VerifyCoinPayments(clientID, secret string, headers http.Header, body []byte) bool {
	if clientID != "" && headers.Get("X-CoinPayments-Client") != clientID {
		return false
	}
	ts := headers.Get("X-CoinPayments-Timestamp")
	sig := headers.Get("X-CoinPayments-Signature")
	if ts == "" || sig == "" {
		return false
	}
	url := "/api/v1/pay/notify/coinpayments"
	expected := coinpaymentsSignature(secret, "POST", url, clientID, ts, body)
	return hmac.Equal([]byte(expected), []byte(sig))
}

func init() { Register(CoinPaymentsGateway{}) }

// ---- PayPal Checkout v2 (orders API) ----

type PayPalGateway struct{}

func (PayPalGateway) Name() string  { return "paypal" }
func (PayPalGateway) Enabled() bool { return payBool("pay.paypal_enabled") }

func (g PayPalGateway) base() string {
	if paySet("pay.paypal_env", "sandbox") == "live" {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

func (g PayPalGateway) token() (string, error) {
	id := paySet("pay.paypal_client_id", "")
	secret := paySet("pay.paypal_secret", "")
	if id == "" || secret == "" {
		return "", fmt.Errorf("paypal 未配置 client_id/secret")
	}
	httpReq, err := http.NewRequest("POST", g.base()+"/v1/oauth2/token",
		strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	httpReq.SetBasicAuth(id, secret)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.AccessToken == "" {
		return "", fmt.Errorf("paypal token 获取失败")
	}
	return out.AccessToken, nil
}

func init() { Register(PayPalGateway{}) }

func (g PayPalGateway) Create(req CreateRequest) (CreateResult, error) {
	token, err := g.token()
	if err != nil {
		return CreateResult{}, err
	}

	amount := fmt.Sprintf("%.2f", float64(req.AmountCents)/100)
	payload := map[string]any{
		"intent": "CAPTURE",
		"purchase_units": []map[string]any{{
			"custom_id": req.OrderNo,
			"amount":    map[string]any{"currency_code": req.Currency, "value": amount},
		}},
		"application_context": map[string]any{
			"return_url": req.ReturnURL,
			"cancel_url": req.ReturnURL,
		},
	}
	body, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequest("POST", g.base()+"/v2/checkout/orders", bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return CreateResult{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out struct {
		ID    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.ID == "" {
		return CreateResult{}, fmt.Errorf("paypal 下单失败: %s", strings.TrimSpace(string(raw)))
	}
	payURL := ""
	for _, l := range out.Links {
		if l.Rel == "approve" {
			payURL = l.Href
		}
	}
	return CreateResult{PayURL: payURL, Raw: map[string]any{"paypal_order_id": out.ID}}, nil
}

// VerifyPayPal confirms a webhook event through PayPal's verify API.
func VerifyPayPal(clientID, secret, env string, headers http.Header, body []byte) bool {
	base := "https://api-m.sandbox.paypal.com"
	if env == "live" {
		base = "https://api-m.paypal.com"
	}
	token, err := PayPalGateway{}.token()
	if err != nil {
		return false
	}
	verifyBody := map[string]any{
		"transmission_id":   headers.Get("Paypal-Transmission-Id"),
		"transmission_time": headers.Get("Paypal-Transmission-Time"),
		"transmission_sig":  headers.Get("Paypal-Transmission-Sig"),
		"cert_url":          headers.Get("Paypal-Cert-Url"),
		"auth_algo":         headers.Get("Paypal-Auth-Algo"),
		"webhook_id":        paySet("pay.paypal_webhook_id", ""),
		"webhook_event":     json.RawMessage(body),
	}
	vb, _ := json.Marshal(verifyBody)
	httpReq, _ := http.NewRequest("POST", base+"/v1/notifications/verify-webhook-signature", bytes.NewReader(vb))
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out struct {
		VerificationStatus string `json:"verification_status"`
	}
	if json.Unmarshal(raw, &out) != nil {
		return false
	}
	return out.VerificationStatus == "SUCCESS"
}
