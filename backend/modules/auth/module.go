package auth

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	"reflexcms/backend/database/migrations"
	authcontrollers "reflexcms/backend/modules/auth/http/controllers"
)

// Module owns credential verification and opaque admin sessions. Part of the
// platform core: every product composition includes it.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "auth" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.AuthMigrations()
}

func (m *Module) Routes() {
	controller := authcontrollers.NewAuthController()
	r := facades.Route()

	r.Post("/api/auth/login", controller.Login)
	r.Get("/api/auth/me", controller.Me)
	r.Post("/api/auth/logout", controller.Logout)

	// Public-frontend identity probe: 200 with {user: null} when anonymous so
	// pages can render logged-out state without treating it as an error.
	r.Get("/api/v1/me", func(ctx http.Context) http.Response {
		if identity, ok := authcontrollers.RequireIdentity(ctx); ok {
			return ctx.Response().Success().Json(http.Json{"user": identity})
		}
		return ctx.Response().Success().Json(http.Json{"user": nil})
	})
}

func (m *Module) Boot() {}
