package config

import (
	"strings"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	// CORS_ORIGINS: comma-separated whitelist for the Phase B public
	// frontend. Empty (default) keeps CORS fully disabled — the admin BFF is
	// same-origin and needs none.
	raw, _ := config.Env("CORS_ORIGINS", "").(string)
	origins := splitOrigins(raw)
	hasOrigins := len(origins) > 0

	config.Add("cors", map[string]any{
		"paths":           corsPaths(hasOrigins),
		"allowed_methods": []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		"allowed_origins": origins,
		"allowed_headers": []string{"Content-Type", "Authorization"},
		"exposed_headers": []string{},
		"max_age":         3600,
		// Credentials travel as the httpOnly admin_session cookie, so the
		// browser only needs them when a whitelisted origin is configured.
		"supports_credentials": hasOrigins,
	})
}

func corsPaths(hasOrigins bool) []string {
	if !hasOrigins {
		return []string{}
	}
	return []string{"api/v1/*"}
}

func splitOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
