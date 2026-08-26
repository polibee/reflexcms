package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260824000001CreateRolesTable struct{}

func (r *M20260824000001CreateRolesTable) Signature() string {
	return "20260824000001_create_roles_table"
}

func (r *M20260824000001CreateRolesTable) Up() error {
	if facades.Schema().HasTable("roles") {
		return nil
	}
	return facades.Schema().Create("roles", func(table schema.Blueprint) {
		table.ID()
		table.String("name", 64)
		table.String("display_name", 128)
		// permissions lives as JSONB; the model layer owns its codec.
		table.Jsonb("permissions")
		table.Integer("sort").Default(100)
		table.TimestampsTz()
		table.Unique("name")
	})
}

func (r *M20260824000001CreateRolesTable) Down() error {
	return facades.Schema().DropIfExists("roles")
}
