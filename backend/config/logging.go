package config

import (
	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("logging", map[string]any{
		"default": config.Env("LOG_CHANNEL", "stack"),
		// Available Drivers: "single", "daily", "custom", "stack"
		"channels": map[string]any{
			"stack": map[string]any{
				"driver":   "stack",
				"channels": []string{"daily"},
			},
			"single": map[string]any{
				"driver":    "single",
				"path":      "storage/logs/goravel.log",
				"level":     config.Env("LOG_LEVEL", "debug"),
				"print":     false,
				"formatter": "text",
			},
			"daily": map[string]any{
				"driver":    "daily",
				"path":      "storage/logs/goravel.log",
				"level":     config.Env("LOG_LEVEL", "debug"),
				"days":      7,
				"print":     false,
				"formatter": "text",
			},
			"otel": map[string]any{
				"driver":          "otel",
				"instrument_name": config.GetString("APP_NAME", "goravel/log"),
			},
		},
	})
}
