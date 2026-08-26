package notification

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	"reflexcms/backend/database/migrations"
	adminhub "reflexcms/backend/modules/adminhub"
	authcontrollers "reflexcms/backend/modules/auth/http/controllers"
	notificationmodels "reflexcms/backend/modules/notification/models"
	notificationservices "reflexcms/backend/modules/notification/services"
)

// Module is the in-app notification feature module. It owns the
// notifications table and the read-all endpoint; triggers are wired by the
// composition root (products/) so this module never imports emitters.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "notification" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.NotificationMigrations()
}

func (m *Module) Routes() {
	adminhub.Register(&adminhub.Spec{
		Name:             "notifications",
		Model:            notificationmodels.Notification{},
		Searchable:       []string{"type"},
		Sortable:         []string{"id", "read_at", "created_at"},
		Fillable:         []string{}, // system-generated: no direct writes
		PermissionPrefix: "notifications",
	})

	facades.Route().Post("/api/admin/notifications/read-all", func(ctx http.Context) http.Response {
		identity, ok := authcontrollers.RequireIdentity(ctx)
		if !ok {
			return httpx.Error(ctx, 401, "Unauthorized")
		}
		rows, err := notificationservices.MarkAllRead(identity.ID)
		if err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"updated": rows})
	})
}

func (m *Module) Boot() {}
