package config

import (
	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("jwt", map[string]any{
		// JWT Authentication Secret.
		// Generate with: go run . artisan jwt:secret
		"secret": config.Env("JWT_SECRET", ""),
		// JWT time to live (minutes). 0 means never expiring (not recommended).
		"ttl": config.Env("JWT_TTL", 60),
		// Refresh time to live (minutes).
		"refresh_ttl": config.Env("JWT_REFRESH_TTL", 20160),
	})
}
