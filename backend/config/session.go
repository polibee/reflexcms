package config

import (
	"github.com/goravel/framework/support/path"
	"github.com/goravel/framework/support/str"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("session", map[string]any{
		// Available Drivers: "file", "custom"
		"default": config.Env("SESSION_DRIVER", "file"),
		"drivers": map[string]any{
			"file": map[string]any{
				"driver": "file",
			},
		},
		"lifetime":        config.Env("SESSION_LIFETIME", 120),
		"expire_on_close": config.Env("SESSION_EXPIRE_ON_CLOSE", false),
		"files":           path.Storage("framework/sessions"),
		"gc_interval":     config.Env("SESSION_GC_INTERVAL", 30),
		"cookie":          config.Env("SESSION_COOKIE", str.Of(config.GetString("app.name")).Snake().Lower().String()+"_session"),
		"path":            config.Env("SESSION_PATH", "/"),
		"domain":          config.Env("SESSION_DOMAIN", ""),
		"secure":          config.Env("SESSION_SECURE", false),
		"http_only":       config.Env("SESSION_HTTP_ONLY", true),
		"same_site":       config.Env("SESSION_SAME_SITE", "lax"),
	})
}
