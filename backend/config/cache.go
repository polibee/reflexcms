package config

import (
	"github.com/goravel/framework/contracts/cache"
	redisfacades "github.com/goravel/redis/facades"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("cache", map[string]any{
		// Default Cache Store
		//
		// Set CACHE_STORE=redis in .env to switch the whole cache layer to Redis.
		"default": config.Env("CACHE_STORE", "memory"),

		// Cache Stores
		//
		// Available Drivers: "memory", "custom"
		"stores": map[string]any{
			"memory": map[string]any{
				"driver": "memory",
			},
			"redis": map[string]any{
				"driver":     "custom",
				"connection": "default",
				"via": func() (cache.Driver, error) {
					return redisfacades.Cache("redis") // The `redis` value is the key of `stores`
				},
			},
		},

		// Cache Key Prefix
		//
		// Must be a-zA-Z0-9_-
		"prefix": config.GetString("APP_NAME", "goravel") + "_cache",
	})
}
