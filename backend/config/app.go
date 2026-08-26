package config

import (
	"github.com/goravel/framework/support/carbon"

	"reflexcms/backend/app/facades"
)

// Boot Start all init methods of the current folder to bootstrap all config.
func Boot() {}

func init() {
	config := facades.Config()
	config.Add("app", map[string]any{
		"name":  config.Env("APP_NAME", "ReflexCMS"),
		"env":   config.Env("APP_ENV", "production"),
		"debug": config.Env("APP_DEBUG", false),
		"timezone": carbon.UTC,
		"locale":          "zh-CN",
		"fallback_locale": "en",
		// Encryption Key: 32 character string. Generate with:
		//   go run . artisan key:generate
		"key": config.Env("APP_KEY", ""),
		"maintenance": map[string]any{
			"driver": config.Env("APP_MAINTENANCE_DRIVER", "file"),
			"store":  config.Env("APP_MAINTENANCE_STORE", ""),
		},
	})
}
