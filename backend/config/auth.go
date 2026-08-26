package config

import (
	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("auth", map[string]any{
		"defaults": map[string]any{
			"guard": "user",
		},
		// Supported drivers: "jwt", "session"
		"guards": map[string]any{
			"user": map[string]any{
				"driver":   "jwt",
				"provider": "user",
			},
		},
		// Supported: "orm"
		"providers": map[string]any{
			"user": map[string]any{
				"driver": "orm",
			},
		},
	})
}
