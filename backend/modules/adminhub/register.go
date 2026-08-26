package adminhub

import (
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
)

// RegisterCoreRoutes wires everything owned by the platform core: the
// generic resource gateway and dashboard aggregates. Feature-module routes
// (auth trio, article actions, topic moderation…) are registered by their
// own modules during Boot.
//
// Resource specs are looked up per-request from the registry, so modules can
// Register(spec) any time before traffic flows.
func RegisterCoreRoutes() {
	gateway := NewGatewayController()
	stats := NewStatsController()

	r := facades.Route()

	// Static routes MUST be registered before wildcard routes so gin
	// prioritizes them in the radix tree.
	r.Get("/api/admin/settings-all", SettingsAll)
	r.Put("/api/admin/settings-batch", SettingsBatch)
	r.Get("/api/admin/stats", stats.Index)

	r.Get("/api/admin/{resource}", gateway.Index)
	r.Post("/api/admin/{resource}", gateway.Store)
	r.Get("/api/admin/{resource}/{id}", gateway.Show)
	r.Put("/api/admin/{resource}/{id}", gateway.Update)
	r.Delete("/api/admin/{resource}/{id}", gateway.Destroy)
	r.Post("/api/admin/{resource}/bulk-delete", gateway.BulkDelete)
	// Recycle bin: only specs with SoftDeletable accept it (others 404).
	r.Post("/api/admin/{resource}/{id}/restore", gateway.Restore)
	r.Post("/api/admin/uploads", Upload)
	r.Static("storage", "./storage/app")
}

// RegisterBulkDeleteFor re-binds the static bulk-delete path for a resource
// whose {id}/{action} subtrees shadow the wildcard bulk-delete route.
func RegisterBulkDeleteFor(resource string) {
	gateway := NewGatewayController()
	facades.Route().Post("/api/admin/"+resource+"/bulk-delete", func(ctx http.Context) http.Response {
		return gateway.BulkDeleteForResource(ctx, resource)
	})
}
