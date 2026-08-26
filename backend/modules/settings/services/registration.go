package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"reflexcms/backend/app/facades"
)

var mainstreamDomains = map[string]bool{
	"gmail.com": true, "outlook.com": true, "hotmail.com": true,
	"yahoo.com": true, "icloud.com": true,
	"qq.com": true, "163.com": true, "126.com": true, "sina.com": true,
	"foxmail.com": true, "proton.me": true, "protonmail.com": true,
}

func emailDomain(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

func parseStringList(raw string) []string {
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func containsFold(list []string, item string) bool {
	for _, v := range list {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}

// ValidateRegistrationEmail checks the registration restriction chain
// (plan ADR-014): blacklist → whitelist → mainstream_only.
func ValidateRegistrationEmail(email string) error {
	domain := emailDomain(email)
	if domain == "" {
		return fmt.Errorf("invalid email address")
	}

	for _, d := range defaultBlacklist {
		if strings.EqualFold(domain, d) {
			return fmt.Errorf("该邮箱域名不可用")
		}
	}

	blacklist := parseStringList(Get(KeyRegEmailBlacklist))
	for _, d := range blacklist {
		if strings.EqualFold(domain, d) {
			return fmt.Errorf("该邮箱域名不可用")
		}
	}

	whitelist := parseStringList(Get(KeyRegEmailWhitelist))
	if len(whitelist) > 0 && !containsFold(whitelist, domain) {
		return fmt.Errorf("该邮箱域名不在白名单内")
	}

	if GetBoolSetting(KeyRegMainstreamOnly) && !mainstreamDomains[domain] {
		return fmt.Errorf("仅支持主流邮箱注册（Gmail、Outlook、QQ 邮箱等）")
	}

	return nil
}

// TurnstileEnabled reports whether Turnstile verification should be enforced.
// Only active in production — development skips captcha entirely.
func TurnstileEnabled() bool {
	return GetBoolSetting(KeyRegTurnstileEnabled) && facades.Config().GetString("app.env", "local") == "production"
}

var defaultBlacklist = []string{
	"tempmail.com", "10minutemail.com", "guerrillamail.com",
	"mailinator.com", "yopmail.com", "throwawaymail.com",
}
