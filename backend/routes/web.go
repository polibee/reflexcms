package routes

import (
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
)

// Web registers operational root routes.
func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().Success().Json(http.Json{
			"name":    facades.Config().GetString("app.name", "ReflexCMS"),
			"service": "reflexcms-backend",
			"docs":    "/healthz",
		})
	})
}
