package config

import (
	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("mail", map[string]any{
		"host": config.Env("MAIL_HOST", ""),
		"port": config.Env("MAIL_PORT", 587),
		"from": map[string]any{
			"address": config.Env("MAIL_FROM_ADDRESS", "hello@example.com"),
			"name":    config.Env("MAIL_FROM_NAME", "Example"),
		},
		"username": config.Env("MAIL_USERNAME"),
		"password": config.Env("MAIL_PASSWORD"),
		// Available Drivers: "html", "custom"
		"template": map[string]any{
			"default": config.Env("MAIL_TEMPLATE_ENGINE", "html"),
			"engines": map[string]any{
				"html": map[string]any{
					"driver": "html",
					"path":   config.Env("MAIL_VIEWS_PATH", "resources/views/mail"),
				},
			},
		},
	})
}
