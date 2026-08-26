// Package kernel defines the module contract every ReflexCMS feature unit
// implements. Products (see products/ + cmd/) compose modules; modules never
// import each other (docs/模块化架构与产品组合.md §2).
package kernel

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/schedule"
	databaseseeder "github.com/goravel/framework/contracts/database/seeder"
)

type Module interface {
	// Name identifies the module, e.g. "cms".
	Name() string
	// Migrations returns exactly the schema changes this module owns.
	Migrations() []schema.Migration
	// Routes registers the module's HTTP routes (and its gateway specs).
	// Called during the framework routing phase — the only point where
	// facades.Route() is guaranteed to be resolvable.
	Routes()
	// Boot wires post-provider concerns that are not route-bound.
	Boot()
}

// SeederProvider is implemented by modules that own database seeders.
type SeederProvider interface {
	Seeders() []databaseseeder.Seeder
}

// EventProvider is implemented by modules that register event listeners.
type EventProvider interface {
	Events() map[event.Event][]event.Listener
}

// SchedulerProvider is implemented by modules that own scheduled tasks.
type SchedulerProvider interface {
	Schedule() []schedule.Event
}
