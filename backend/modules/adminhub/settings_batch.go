package adminhub

import (
	"encoding/json"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	settingsservices "reflexcms/backend/modules/settings/services"
)

// SettingsAll handles GET /api/admin/settings-all — returns all settings
// grouped by their group field as { group: { key: decodedValue } }.
func SettingsAll(ctx http.Context) http.Response {
	if resp := CheckPermission(ctx, "settings.view"); resp != nil {
		return *resp
	}

	var rows []map[string]any
	if err := facades.Orm().Query().Table("settings").
		Select("key", "value", `"group"`).
		OrderBy(`"group"`, "asc").OrderBy("key", "asc").
		Get(&rows); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	result := map[string]map[string]any{}
	for _, row := range rows {
		group, _ := row["group"].(string)
		key, _ := row["key"].(string)
		val := row["value"]

		if result[group] == nil {
			result[group] = map[string]any{}
		}
		// Decode JSON values for frontend consumption.
		if s, ok := val.(string); ok && s != "" {
			var decoded any
			if json.Unmarshal([]byte(s), &decoded) == nil {
				val = decoded
			}
		}
		result[group][key] = val
	}

	return ctx.Response().Success().Json(result)
}

// SettingsBatch handles PUT /api/admin/settings-batch — accepts
// { "values": { "site.mode": "forum", "reg.enabled": true } } and upserts.
func SettingsBatch(ctx http.Context) http.Response {
	if resp := CheckPermission(ctx, "settings.edit"); resp != nil {
		return *resp
	}

	var req struct {
		Values map[string]any `json:"values"`
	}
	if err := ctx.Request().Bind(&req); err != nil || len(req.Values) == 0 {
		return httpx.Error(ctx, 422, "values object is required")
	}

	saved := 0
	for key, value := range req.Values {
		if err := settingsservices.ValidateSettingValue(key, value); err != nil {
			return httpx.Error(ctx, 422, err.Error())
		}
		group := "general"
		if idx := indexOfDot(key); idx > 0 {
			group = key[:idx]
		}
		if err := settingsservices.SetWithGroup(key, value, group); err != nil {
			return httpx.Error(ctx, 500, "failed to save "+key+": "+err.Error())
		}
		facades.Cache().Forget("settings:" + key)
		saved++
	}

	return ctx.Response().Success().Json(http.Json{"saved": saved})
}

func indexOfDot(s string) int {
	for i, c := range s {
		if c == '.' {
			return i
		}
	}
	return -1
}
