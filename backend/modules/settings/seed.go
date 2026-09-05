package settings

import (
	settingsservices "reflexcms/backend/modules/settings/services"
)

type seedEntry struct {
	key   string
	value any
	group string
}

func seedDefaults() {
	defaults := []seedEntry{
		{"site.mode", "hybrid", "site"},
		{"site.home", "both", "site"},
		{"site.locale", "zh-CN", "site"},
		{"site.name", "ReflexCMS", "site"},
		{"site.description", "", "site"},
		{"site.icp", "", "site"},
		{"site.keywords", "", "site"},
		{"reg.enabled", true, "registration"},
		{"reg.invite_required", false, "registration"},
		{"reg.email_whitelist", []string{}, "registration"},
		{"reg.email_blacklist", []string{}, "registration"},
		{"reg.mainstream_only", false, "registration"},
		{"reg.turnstile_enabled", false, "registration"},
		{"reg.turnstile_site_key", "", "registration"},
		{"reg.turnstile_secret_key", "", "registration"},
		{"email.driver", "log", "email"},
		{"email.from_address", "", "email"},
		{"email.from_name", "ReflexCMS", "email"},
		{"email.resend_api_key", "", "email"},
		{"email.smtp_host", "", "email"},
		{"email.smtp_port", 587, "email"},
		{"email.smtp_user", "", "email"},
		{"email.smtp_pass", "", "email"},
		{"seo.google_site_verification", "", "seo"},
		{"seo.baidu_verification", "", "seo"},
		{"seo.ads_txt", "", "seo"},
		{"seo.head_scripts", "", "seo"},
		{"analytics.ga_property_id", "", "analytics"},
		{"analytics.ga_credentials_json", "", "analytics"},
		{"ai.provider", "", "ai"},
		{"ai.model", "", "ai"},
		{"ai.api_key", "", "ai"},
		{"ai.base_url", "", "ai"},
		{"privacy.cookie_banner", true, "privacy"},
		{"privacy.privacy_policy_url", "", "privacy"},
		{"privacy.policy_updated_at", "", "privacy"},
		{"points.currency_name", "鸡腿", "points"},
		{"points.earn_topic", 5, "points"},
		{"points.topic_daily_cap", 20, "points"},
		{"points.earn_reply", 2, "points"},
		{"points.reply_daily_cap", 20, "points"},
		{"points.earn_signin", 20, "points"},
		{"points.level_thresholds", "0,100,300,600,1000,1500,2200,3000", "points"},
		{"points.level_names", "新手,活跃,达人,熟练,精英,元老,传奇", "points"},
		{"points.level_perks", "发帖与回复|发布主题帖|自定义个人卡片|专属徽章|优先客服支持|专属版块访问|限量邀请码折扣", "points"},
		{"pay.shop_enabled", false, "pay"},
		{"pay.display_currency", "USD", "pay"},
		{"pay.aff_xcash", "https://dash.xca.sh/register?ref=2GWV5MKT", "pay"},
		{"pay.aff_coinpayments", "", "pay"},
		{"pay.aff_nowpayments", "https://account.nowpayments.io/create-account?link_id=3940543227", "pay"},
		{"pay.aff_xunhupay", "", "pay"},
		{"pay.aff_codepay", "", "pay"},
		{"pay.aff_paypal", "", "pay"},
		{"pay.nowpayments_enabled", false, "pay"},
		{"pay.nowpayments_api_key", "", "pay"},
		{"pay.xcash_enabled", false, "pay"},
		{"pay.xcash_appid", "", "pay"},
		{"pay.xcash_hmac_key", "", "pay"},
		{"pay.xcash_api", "https://pay.xca.sh", "pay"},
		{"pay.coinpayments_enabled", false, "pay"},
		{"pay.coinpayments_client_id", "", "pay"},
		{"pay.coinpayments_client_secret", "", "pay"},
		{"pay.xunhu_enabled", false, "pay"},
		{"pay.xunhu_appid", "", "pay"},
		{"pay.xunhu_appsecret", "", "pay"},
		{"pay.codepay_enabled", false, "pay"},
		{"pay.codepay_id", "", "pay"},
		{"pay.codepay_key", "", "pay"},
		{"pay.paypal_enabled", false, "pay"},
		{"pay.paypal_client_id", "", "pay"},
		{"pay.paypal_secret", "", "pay"},
		{"pay.paypal_env", "sandbox", "pay"},
		{"pay.paypal_webhook_id", "", "pay"},
	}

	for _, d := range defaults {
		if settingsservices.Get(d.key) != "" {
			continue
		}
		_ = settingsservices.SetWithGroup(d.key, d.value, d.group)
	}
}
