package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000003CreateRepliesTable struct{}

func (r *M20260825000003CreateRepliesTable) Signature() string {
	return "20260825000003_create_replies_table"
}

func (r *M20260825000003CreateRepliesTable) Up() error {
	if facades.Schema().HasTable("replies") {
		return nil
	}
	return facades.Schema().Create("replies", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("topic_id")
		table.UnsignedBigInteger("user_id").Nullable()
		table.UnsignedBigInteger("parent_id").Nullable()
		table.Text("content")
		table.Integer("floor")
		table.UnsignedBigInteger("like_count").Default(0)
		table.SoftDeletes()
		table.TimestampsTz()
		table.Index("topic_id")
		table.Index("user_id")
		table.Index("parent_id")
	})
}

func (r *M20260825000003CreateRepliesTable) Down() error {
	return facades.Schema().DropIfExists("replies")
}
