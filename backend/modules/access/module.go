package access

import (
	"github.com/goravel/framework/contracts/database/schema"
	databaseseeder "github.com/goravel/framework/contracts/database/seeder"

	"reflexcms/backend/database/migrations"
	"reflexcms/backend/database/seeders"
)

// Module is the RBAC platform module (roles, users). It is part of the
// platform core: every product composition includes it.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "access" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.AccessMigrations()
}

func (m *Module) Seeders() []databaseseeder.Seeder {
	return seeders.AccessSeeders()
}

func (m *Module) Routes() {
	Register()
}

func (m *Module) Boot() {}
