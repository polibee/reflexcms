package settings

import (
	"encoding/json"
	"strings"

	"reflexcms/backend/app/facades"

	adminhub "reflexcms/backend/modules/adminhub"
	settingsmodels "reflexcms/backend/modules/settings/models"
	settingsservices "reflexcms/backend/modules/settings/services"
)

// Site mode enum (plan §5.4). Runtime-switchable, zero migrations.
const (
	ModeBlog   = "blog"
	ModeForum  = "forum"
	ModeHybrid = "hybrid"
)

func isValidMode(v string) bool {
	return v == ModeBlog || v == ModeForum || v == ModeHybrid
}

func settingSpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "settings",
		Model:      settingsmodels.Setting{},
		Searchable: []string{"key"},
		Sortable:   []string{"id", "key", "group", "updated_at"},
		Fillable:   []string{"key", "value", "group"},
		Transform:  decodeSettingValue,
		BeforeSave: prepareSetting,
		AfterSave: func(instance any, data map[string]any) error {
			if s, ok := instance.(settingsmodels.Setting); ok && s.Key != "" {
				facades.Cache().Forget("settings:" + s.Key)
			}
			return nil
		},
		PermissionPrefix: "settings",
	}
}

// decodeSettingValue unwraps the stored JSON encoding so admin forms bind
// plain values ("forum", not "\"forum\"") — enabling visual select inputs.
func decodeSettingValue(row any) map[string]any {
	b, err := json.Marshal(row)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{}
	}
	if raw, ok := m["value"].(string); ok && raw != "" {
		var decoded any
		if err := json.Unmarshal([]byte(raw), &decoded); err == nil {
			m["value"] = decoded
		}
	}
	return m
}

// prepareSetting re-encodes the submitted value into its JSON storage form
// and enforces the site.mode enum. Values arrive decoded from visual forms;
// pre-encoded JSON strings pass through untouched.
func prepareSetting(data map[string]any) error {
	raw, present := data["value"]
	if !present {
		return nil
	}

	if key, _ := data["key"].(string); key == "site.mode" {
		mode, _ := raw.(string)
		mode = strings.Trim(mode, `"`)
		if err := settingsservices.ValidateSettingValue("site.mode", mode); err != nil {
			return err
		}
		data["value"] = mode
		return nil
	}

	if key, _ := data["key"].(string); key == "site.home" {
		home, _ := raw.(string)
		home = strings.Trim(home, `"`)
		if err := settingsservices.ValidateSettingValue("site.home", home); err != nil {
			return err
		}
		data["value"] = home
		return nil
	}

	text, isText := raw.(string)
	if !isText {
		return nil // non-string payloads (objects/numbers) store as-is
	}
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		data["value"] = nil
		return nil
	}
	if !json.Valid([]byte(trimmed)) {
		encoded, _ := json.Marshal(text)
		data["value"] = string(encoded)
	} else {
		data["value"] = trimmed
	}
	return nil
}
