package config

import (
	"github.com/goravel/framework/support/path"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("filesystems", map[string]any{
		"default": "local",
		// Supported Drivers: "local", "custom"
		"disks": map[string]any{
			"local": map[string]any{
				"driver": "local",
				"root":   path.Storage("app"),
			},
			"public": map[string]any{
				"driver": "local",
				"root":   path.Storage("app/public"),
				"url":    config.Env("APP_URL", "").(string) + "/storage",
			},
		},
	})
}
