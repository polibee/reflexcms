package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000007CreateNotificationsTable struct{}

func (r *M20260825000007CreateNotificationsTable) Signature() string {
	return "20260825000007_create_notifications_table"
}

func (r *M20260825000007CreateNotificationsTable) Up() error {
	if facades.Schema().HasTable("notifications") {
		return nil
	}
	return facades.Schema().Create("notifications", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.String("type", 64)
		// JSON payload; the notification service owns the codec.
		table.Jsonb("data")
		table.DateTimeTz("read_at").Nullable()
		table.TimestampsTz()
		table.Index("user_id")
	})
}

func (r *M20260825000007CreateNotificationsTable) Down() error {
	return facades.Schema().DropIfExists("notifications")
}
