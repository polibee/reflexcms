package services

// Setting keys grouped by domain (docs/设置系统与插件架构.md §1).
const (
	KeySiteMode        = "site.mode"
	KeySiteName        = "site.name"
	KeySiteDescription = "site.description"
	KeySiteICP         = "site.icp"
	KeySiteKeywords    = "site.keywords"

	KeyRegEnabled          = "reg.enabled"
	KeyRegEmailWhitelist   = "reg.email_whitelist"
	KeyRegEmailBlacklist   = "reg.email_blacklist"
	KeyRegMainstreamOnly   = "reg.mainstream_only"
	KeyRegTurnstileEnabled = "reg.turnstile_enabled"
	KeyRegTurnstileSiteKey = "reg.turnstile_site_key"
	KeyRegTurnstileSecret  = "reg.turnstile_secret_key"

	KeyEmailDriver      = "email.driver"
	KeyEmailFromAddress = "email.from_address"
	KeyEmailFromName    = "email.from_name"
	KeyEmailResendKey   = "email.resend_api_key"

	KeySEOGoogleVerification = "seo.google_site_verification"
	KeySEOAdsTxt             = "seo.ads_txt"
	KeySEOHeadScripts        = "seo.head_scripts"

	KeyAIBaseURL = "ai.base_url"

	KeyPrivacyCookieBanner = "privacy.cookie_banner"
)
