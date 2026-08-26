package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260826000003CreateMessagesTable struct{}

func (r *M20260826000003CreateMessagesTable) Signature() string {
	return "20260826000003_create_messages_table"
}

func (r *M20260826000003CreateMessagesTable) Up() error {
	if facades.Schema().HasTable("messages") {
		return nil
	}
	return facades.Schema().Create("messages", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("sender_id").Nullable()
		table.UnsignedBigInteger("recipient_id")
		table.Text("body")
		table.DateTimeTz("read_at").Nullable()
		table.TimestampsTz()
		table.Index("recipient_id")
		table.Index("sender_id")
	})
}

func (r *M20260826000003CreateMessagesTable) Down() error {
	return facades.Schema().Drop("messages")
}
