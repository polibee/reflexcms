package bootstrap

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/schema"
	databaseseeder "github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/event"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/foundation"

	"reflexcms/backend/app/console/commands"
	"reflexcms/backend/config"
	"reflexcms/backend/database/migrations"
	"reflexcms/backend/kernel"
	adminhub "reflexcms/backend/modules/adminhub"
	"reflexcms/backend/routes"
)

// Boot composes the application from the given modules. The composition —
// which modules run, i.e. which product this binary is — lives entirely with
// the caller (root main.go = Full; cmd/blog, cmd/forum = single-product).
func Boot(mods ...kernel.Module) contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(func() []schema.Migration {
			all := migrations.PlatformMigrations()
			for _, m := range mods {
				all = append(all, m.Migrations()...)
			}
			return all
		}).
		WithSeeders(func() []databaseseeder.Seeder {
			var all []databaseseeder.Seeder
			for _, m := range mods {
				if p, ok := m.(kernel.SeederProvider); ok {
					all = append(all, p.Seeders()...)
				}
			}
			return all
		}).
		WithEvents(func() map[event.Event][]event.Listener {
			registered := map[event.Event][]event.Listener{}
			for _, m := range mods {
				if p, ok := m.(kernel.EventProvider); ok {
					for evt, listeners := range p.Events() {
						registered[evt] = append(registered[evt], listeners...)
					}
				}
			}
			return registered
		}).
		WithSchedule(func() []schedule.Event {
			var events []schedule.Event
			for _, m := range mods {
				if p, ok := m.(kernel.SchedulerProvider); ok {
					events = append(events, p.Schedule()...)
				}
			}
			return events
		}).
		WithRouting(func() {
			routes.Web()
			routes.Api()
			// Platform-core owned routes, then each feature module's own.
			adminhub.RegisterCoreRoutes()
			for _, m := range mods {
				m.Routes()
			}
		}).
		WithProviders(func() []contractsfoundation.ServiceProvider {
			// Framework providers first (route, orm, cache…), then the
			// composition driver that Boots every module.
			return append(Providers(), NewModuleProviders(mods))
		}).
		WithCommands(func() []console.Command {
			return []console.Command{
				&commands.AdminBootstrap{},
			}
		}).
		WithConfig(config.Boot).
		Create()
}

// NewModuleProviders wraps the composition into a framework service provider.
func NewModuleProviders(mods []kernel.Module) contractsfoundation.ServiceProvider {
	return &moduleServiceProvider{mods: mods}
}
