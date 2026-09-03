// Legacy-MD5 payment gateways (虎皮椒 XunhuPay, 码支付 CodePay).
//
// PROTOCOL/POLICY CONFLICT: these providers sign parameters with MD5 — a
// primitive this project's security policy (Mimosa) forbids. The gateways
// are therefore wired but DISABLED until the policy owner explicitly
// allows the provider-mandated digest: Create() returns a configuration
// error pointing at this decision, and the callback verifiers stay
// unimplemented (they would need the same primitive). Nothing in this file
// touches the weak primitive; enabling these gateways is a deliberate,
// documented policy exception, not an accident.
package gateways

import (
	"fmt"

	settingsservices "reflexcms/backend/modules/settings/services"
)

// md5PolicyBlocked is the sentinel error for gateways whose wire protocol
// requires a primitive the security policy disallows.
var md5PolicyBlocked = fmt.Errorf(
	"该渠道协议要求 MD5 签名，与本项目安全策略冲突：需策略所有者显式允许后才能启用（见 gateways/legacy_md5_gateways.go 文件头说明）",
)

func legacyProviderEnabled(key string) bool {
	if payBool(key) {
		settingsservices.Get(key) // presence documents intent in settings
	}
	return false
}

// ---- 虎皮椒 XunhuPay ----

type XunhuGateway struct{}

func (XunhuGateway) Name() string    { return "xunhupay" }
func (XunhuGateway) Enabled() bool   { return legacyProviderEnabled("pay.xunhu_enabled") }

func (XunhuGateway) Create(req CreateRequest) (CreateResult, error) {
	// Wire format (per provider docs): appid, trade_order_id, total_fee,
	// title, time, notify_url, return_url, nonce_str, hash — where hash is
	// MD5-of-UPPER(MD5(secret + sorted "k=v" joined by &)). Blocked by
	// security policy; see file header.
	_ = req
	return CreateResult{}, md5PolicyBlocked
}

func init() { Register(XunhuGateway{}) }

// ---- 码支付 CodePay ----

type CodePayGateway struct{}

func (CodePayGateway) Name() string    { return "codepay" }
func (CodePayGateway) Enabled() bool   { return legacyProviderEnabled("pay.codepay_enabled") }

func (CodePayGateway) Create(req CreateRequest) (CreateResult, error) {
	// Wire format (per provider docs): GET creat_order/ with id, type,
	// out_trade_no, price, notify_url, return_url, name, sign — sign is
	// MD5(sorted "k=v" joined by & + KEY). Blocked by security policy; see
	// file header.
	_ = req
	return CreateResult{}, md5PolicyBlocked
}

func init() { Register(CodePayGateway{}) }
