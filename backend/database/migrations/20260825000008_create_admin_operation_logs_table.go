package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000008CreateAdminOperationLogsTable struct{}

func (r *M20260825000008CreateAdminOperationLogsTable) Signature() string {
	return "20260825000008_create_admin_operation_logs_table"
}

func (r *M20260825000008CreateAdminOperationLogsTable) Up() error {
	if facades.Schema().HasTable("admin_operation_logs") {
		return nil
	}
	return facades.Schema().Create("admin_operation_logs", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id").Nullable()
		table.String("action", 32)
		table.String("resource", 64).Nullable()
		table.String("resource_id", 32).Nullable()
		table.Integer("status_code")
		table.String("ip", 64).Nullable()
		table.String("user_agent", 255).Nullable()
		// Sanitised request snapshot (password/token keys stripped).
		table.Jsonb("payload").Nullable()
		table.DateTimeTz("created_at").UseCurrent()
		table.Index("user_id")
		table.Index("resource")
		table.Index("created_at")
	})
}

func (r *M20260825000008CreateAdminOperationLogsTable) Down() error {
	return facades.Schema().DropIfExists("admin_operation_logs")
}
