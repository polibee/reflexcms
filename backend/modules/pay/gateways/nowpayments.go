package gateways

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"strings"
	"time"
)

// NOWPayments: invoice API with a simple x-api-key header (no request
// signing). Webhooks are verified via the callback signature HMAC if a
// webhook key is configured; otherwise the IPN callback token matches.
type NowPaymentsGateway struct{}

func (NowPaymentsGateway) Name() string  { return "nowpayments" }
func (NowPaymentsGateway) Enabled() bool { return payBool("pay.nowpayments_enabled") }

func (NowPaymentsGateway) Create(req CreateRequest) (CreateResult, error) {
	key := paySet("pay.nowpayments_api_key", "")
	if key == "" {
		return CreateResult{}, fmt.Errorf("nowpayments 未配置 api_key")
	}

	amount := fmt.Sprintf("%.2f", float64(req.AmountCents)/100)
	payload := map[string]any{
		"price_amount":      parseFLoat(amount),
		"price_currency":    strings.ToLower(req.Currency),
		"order_id":          req.OrderNo,
		"order_description": req.Title,
		"ipn_callback_url":  req.NotifyURL,
		"success_url":       req.ReturnURL,
		"cancel_url":        req.ReturnURL,
	}
	body, _ := json.Marshal(payload)

	httpReq, _ := nethttp.NewRequest("POST", "https://api.nowpayments.io/v1/invoice", strings.NewReader(string(body)))
	httpReq.Header.Set("x-api-key", key)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &nethttp.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return CreateResult{}, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var out struct {
		ID      string `json:"id"`
		Invoice string `json:"invoice_url"`
	}
	if err := json.Unmarshal(raw, &out); err != nil || out.Invoice == "" {
		return CreateResult{}, fmt.Errorf("nowpayments 创建失败: %s", strings.TrimSpace(string(raw)))
	}
	return CreateResult{PayURL: out.Invoice, Raw: map[string]any{"invoice_id": out.ID}}, nil
}

func parseFLoat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func init() { Register(NowPaymentsGateway{}) }
