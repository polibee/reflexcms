package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260824000002CreateUsersTable struct{}

func (r *M20260824000002CreateUsersTable) Signature() string {
	return "20260824000002_create_users_table"
}

func (r *M20260824000002CreateUsersTable) Up() error {
	if facades.Schema().HasTable("users") {
		return nil
	}
	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.ID()
		table.String("username", 64)
		table.String("email", 255)
		table.String("password", 255)
		table.String("avatar", 512).Nullable()
		table.Text("bio").Nullable()
		table.String("website", 255).Nullable()
		// Referential integrity is enforced at the application layer for now;
		// indexes keep the hot lookups cheap.
		table.UnsignedBigInteger("role_id").Nullable()
		table.String("status", 32).Default("active")
		table.DateTimeTz("banned_at").Nullable()
		table.SoftDeletes()
		table.TimestampsTz()
		table.Unique("username")
		table.Unique("email")
		table.Index("role_id")
	})
}

func (r *M20260824000002CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
