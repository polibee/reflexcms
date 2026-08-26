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
		{"site.name", "ReflexCMS", "site"},
		{"site.description", "", "site"},
		{"site.icp", "", "site"},
		{"site.keywords", "", "site"},
		{"reg.enabled", true, "registration"},
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
	}

	for _, d := range defaults {
		if settingsservices.Get(d.key) != "" {
			continue
		}
		_ = settingsservices.SetWithGroup(d.key, d.value, d.group)
	}
}
