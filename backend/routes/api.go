package routes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
)

// Api registers operational and versioned public API routes.
// Public content endpoints (Phase B) live under /api/v1 and are gated by site mode.
func Api() {
	facades.Route().Get("/healthz", func(ctx http.Context) http.Response {
		checks := map[string]string{
			"database": "ok",
			"cache":    "ok",
		}

		safeCheck(checks, "database", func() error {
			sqlDB, err := facades.Orm().DB()
			if err != nil {
				return err
			}
			pingCtx, cancel := context.WithTimeout(ctx.Context(), 2*time.Second)
			defer cancel()
			return sqlDB.PingContext(pingCtx)
		})

		safeCheck(checks, "cache", func() error {
			const probeKey = "healthz:probe"
			// Driver factories (e.g. goravel/redis) panic inside the container
			// when their backend is unreachable; safeCheck converts that into a
			// structured degraded result instead of a bare 500.
			if err := facades.Cache().Put(probeKey, time.Now().UnixNano(), time.Minute); err != nil {
				return err
			}
			if facades.Cache().Get(probeKey) == nil {
				return errors.New("probe readback failed")
			}
			facades.Cache().Forget(probeKey)
			return nil
		})

		status := "ok"
		for _, result := range checks {
			if result != "ok" {
				status = "degraded"
				break
			}
		}

		code := http.StatusOK
		if status != "ok" {
			code = http.StatusServiceUnavailable
		}

		return ctx.Response().Json(code, http.Json{
			"status": status,
			"app":    facades.Config().GetString("app.name", "ReflexCMS"),
			"env":    facades.Config().GetString("app.env", "production"),
			"checks": checks,
		})
	})
}

// safeCheck runs one health probe, capturing both returned errors and panics.
func safeCheck(checks map[string]string, name string, fn func() error) {
	defer func() {
		if r := recover(); r != nil {
			checks[name] = fmt.Sprintf("error: %v", r)
		}
	}()

	if err := fn(); err != nil {
		checks[name] = "error: " + err.Error()
	}
}
