package services

import "errors"

/* Enum validation for sensitive setting keys. Lives in the services package
 * so both write paths — the gateway's BeforeSave and the adminhub batch
 * endpoint — enforce identical rules without an import cycle. */

// ValidateSettingValue rejects invalid enum payloads for keys that have
// constrained domains. Keys without rules pass through.
func ValidateSettingValue(key string, value any) error {
	s, _ := value.(string)
	switch key {
	case "site.mode":
		if s != "blog" && s != "forum" && s != "hybrid" {
			return errors.New("site.mode must be blog, forum or hybrid")
		}
	case "site.home":
		if s != "cms" && s != "forum" && s != "both" {
			return errors.New("site.home must be cms, forum or both")
		}
	}
	return nil
}
