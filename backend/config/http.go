package config

import (
	"github.com/goravel/framework/contracts/route"
	ginfacades "github.com/goravel/gin/facades"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		"default": "gin",
		// HTTP Drivers
		"drivers": map[string]any{
			"gin": map[string]any{
				// 1 MiB: must exceed the largest allowed topic/reply body
				// (50000 chars of CJK ≈ 150 KB) with headroom for headers.
				"body_limit":   1024 * 1024,
				"header_limit": 4096,
				"route": func() (route.Route, error) {
					return ginfacades.Route("gin"), nil
				},
			},
		},
		// HTTP URL
		"url": config.Env("APP_URL", "http://localhost"),
		// HTTP Host
		"host": config.Env("APP_HOST", "127.0.0.1"),
		// HTTP Port
		"port": config.Env("APP_PORT", "3000"),
		// HTTP Timeout, default is 3 seconds
		"request_timeout": 3,
		// HTTPS Configuration
		"tls": map[string]any{
			"host": config.Env("APP_HOST", "127.0.0.1"),
			"port": config.Env("APP_PORT", "3000"),
			"ssl": map[string]any{
				"cert": "",
				"key":  "",
			},
		},
		// Default Client Name
		//
		// Outbound requests must go through support/safefetch wrappers built on
		// facades.Http(); never construct raw http.Client elsewhere.
		"default_client": config.Env("HTTP_CLIENT_DEFAULT", "default"),
		"clients": map[string]any{
			"default": map[string]any{
				"base_url":                config.Env("HTTP_CLIENT_BASE_URL", ""),
				"timeout":                 config.Env("HTTP_CLIENT_TIMEOUT", "30s"),
				"max_idle_conns":          config.Env("HTTP_CLIENT_MAX_IDLE_CONNS", 100),
				"max_idle_conns_per_host": config.Env("HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST", 2),
				"max_conns_per_host":      config.Env("HTTP_CLIENT_MAX_CONN_PER_HOST", 0),
				"idle_conn_timeout":       config.Env("HTTP_CLIENT_IDLE_CONN_TIMEOUT", "90s"),
			},
		},
	})
}
