package services

import (
	"encoding/json"
	"strings"

	"reflexcms/backend/app/facades"
	settingsmodels "reflexcms/backend/modules/settings/models"
)

// Site modes (plan §5.4 / O10). The mode is a runtime setting: switching it
// in the admin panel changes the public surface immediately, with zero
// migrations and zero redeployment.
const (
	ModeBlog   = "blog"
	ModeForum  = "forum"
	ModeHybrid = "hybrid"
)

func cacheKey(key string) string { return "settings:" + key }

// Get returns the decoded string value of a setting.
func Get(key string) string {
	raw := GetRaw(key)
	var decoded any
	if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
		switch v := decoded.(type) {
		case string:
			return v
		default:
			b, _ := json.Marshal(decoded)
			return string(b)
		}
	}
	return raw
}

// GetRaw returns the raw JSON string stored under key ("" when unset).
func GetRaw(key string) string {
	raw := facades.Cache().Get(cacheKey(key))
	if s, ok := raw.(string); ok && s != "" {
		return s
	}
	var setting settingsmodels.Setting
	if err := facades.Orm().Query().Where("key = ?", key).FirstOrFail(&setting); err != nil {
		return ""
	}
	_ = facades.Cache().Put(cacheKey(key), setting.Value, 0)
	return setting.Value
}

// GetBoolSetting returns the boolean value of a setting.
func GetBoolSetting(key string) bool {
	v := strings.ToLower(strings.Trim(Get(key), `"`))
	return v == "true" || v == "1"
}

// SetWithGroup stores a value with an explicit group assignment.
func SetWithGroup(key string, value any, group string) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var setting settingsmodels.Setting
	findErr := facades.Orm().Query().Where("key = ?", key).FirstOrFail(&setting)
	if findErr == nil {
		setting.Value = string(encoded)
		setting.Group = group
		return facades.Orm().Query().Save(&setting)
	}

	setting = settingsmodels.Setting{Key: key, Value: string(encoded), Group: group}
	return facades.Orm().Query().Create(&setting)
}

func groupOf(key string) string {
	if strings.HasPrefix(key, "site.") {
		return "site"
	}
	return "general"
}

/* ---- typed helpers for well-known keys ---- */

// Mode returns the active site mode; unknown values degrade to hybrid.
func Mode() string {
	mode := decodeString(Get("site.mode"))
	switch mode {
	case ModeBlog, ModeForum:
		return mode
	}
	return ModeHybrid
}

// SectionEnabled reports whether a public content section is exposed under
// the current site mode. Phase B's /api/v1 middleware is built on this.
func SectionEnabled(section string) bool {
	switch Mode() {
	case ModeBlog:
		return section == "articles" || section == "comments"
	case ModeForum:
		return section == "forums" || section == "topics" || section == "replies"
	default: // hybrid exposes everything
		return true
	}
}

func decodeString(raw string) string {
	var out string
	if err := json.Unmarshal([]byte(raw), &out); err == nil {
		return out
	}
	return strings.Trim(raw, `"`)
}
