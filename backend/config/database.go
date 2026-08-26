package config

import (
	"github.com/goravel/framework/contracts/database/driver"
	postgresfacades "github.com/goravel/postgres/facades"

	"reflexcms/backend/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("database", map[string]any{
		// Default database connection name
		"default": config.Env("DB_CONNECTION"),
		// Database connections
		"connections": map[string]any{
			"postgres": map[string]any{
				"host":     config.Env("DB_HOST"),
				"port":     config.Env("DB_PORT"),
				"database": config.Env("DB_DATABASE"),
				"username": config.Env("DB_USERNAME"),
				"password": config.Env("DB_PASSWORD"),
				"sslmode":  "disable",
				"singular": false,
				"prefix":   "",
				"schema":   config.Env("DB_SCHEMA", "public"),
				"via": func() (driver.Driver, error) {
					return postgresfacades.Postgres("postgres")
				},
			},
		},
		// Redis connections, consumed by the goravel/redis drivers
		// (cache / queue / session).
		"redis": map[string]any{
			"default": map[string]any{
				"host":     config.Env("REDIS_HOST", "127.0.0.1"),
				"password": config.Env("REDIS_PASSWORD", ""),
				"port":     config.Env("REDIS_PORT", 6379),
				"database": config.Env("REDIS_DB", 0),
			},
		},
		// Pool configuration
		"pool": map[string]any{
			"max_idle_conns":     10,
			"max_open_conns":     100,
			"conn_max_idletime":  3600,
			"conn_max_lifetime":  3600,
		},
		// Slow query threshold in milliseconds; slower queries get logged.
		"slow_threshold": 200,
		// Migration Repository Table
		"migrations": map[string]any{
			"table": "migrations",
		},
	})
}
