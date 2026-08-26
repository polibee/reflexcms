package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260824000003CreateAdminSessionsTable struct{}

func (r *M20260824000003CreateAdminSessionsTable) Signature() string {
	return "20260824000003_create_admin_sessions_table"
}

func (r *M20260824000003CreateAdminSessionsTable) Up() error {
	if facades.Schema().HasTable("admin_sessions") {
		return nil
	}
	return facades.Schema().Create("admin_sessions", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		// Only the SHA-256 hex of the opaque token is ever stored.
		table.String("token_hash", 64)
		table.DateTimeTz("expires_at")
		table.String("ip", 64).Nullable()
		table.String("user_agent", 255).Nullable()
		table.DateTimeTz("last_used_at").Nullable()
		table.DateTimeTz("created_at").UseCurrent()
		table.Unique("token_hash")
		table.Index("user_id")
	})
}

func (r *M20260824000003CreateAdminSessionsTable) Down() error {
	return facades.Schema().DropIfExists("admin_sessions")
}
