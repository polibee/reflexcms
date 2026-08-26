package config

import (
	"github.com/goravel/framework/contracts/queue"
	redisfacades "github.com/goravel/redis/facades"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("queue", map[string]any{
		// Default Queue Connection Name
		//
		// Set QUEUE_CONNECTION=redis (or database) for durable queues.
		"default": config.Env("QUEUE_CONNECTION", "sync"),

		// Queue Connections
		//
		// Drivers: "sync", "database", "custom"
		"connections": map[string]any{
			"sync": map[string]any{
				"driver": "sync",
			},
			"database": map[string]any{
				"driver":     "database",
				"connection": config.Env("DB_CONNECTION"),
				"queue":      "default",
				"concurrent": 1,
			},
			"redis": map[string]any{
				"driver":     "custom",
				"connection": "default",
				"queue":      "default",
				"via": func() (queue.Driver, error) {
					return redisfacades.Queue("redis") // The `redis` value is the key of `connections`
				},
			},
		},

		// Failed Queue Jobs
		"failed": map[string]any{
			"database": config.Env("DB_CONNECTION"),
			"table":    "failed_jobs",
		},
	})
}
