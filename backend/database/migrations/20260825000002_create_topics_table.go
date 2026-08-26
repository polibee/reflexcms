package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000002CreateTopicsTable struct{}

func (r *M20260825000002CreateTopicsTable) Signature() string {
	return "20260825000002_create_topics_table"
}

func (r *M20260825000002CreateTopicsTable) Up() error {
	if facades.Schema().HasTable("topics") {
		return nil
	}
	return facades.Schema().Create("topics", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id").Nullable()
		table.UnsignedBigInteger("forum_category_id").Nullable()
		table.String("title", 512)
		table.LongText("content").Nullable()
		table.String("status", 32).Default("open")
		table.Boolean("is_pinned").Default(false)
		table.Boolean("is_featured").Default(false)
		// Denormalised counters: reply_count is maintained inside the reply
		// transaction, like_count is recounted on every toggle.
		table.UnsignedBigInteger("reply_count").Default(0)
		table.UnsignedBigInteger("like_count").Default(0)
		table.UnsignedBigInteger("view_count").Default(0)
		table.DateTimeTz("last_reply_at").Nullable()
		table.UnsignedBigInteger("last_reply_user_id").Nullable()
		table.SoftDeletes()
		table.TimestampsTz()
		table.Index("forum_category_id")
		table.Index("user_id")
		table.Index("status")
		table.Index("last_reply_at")
	})
}

func (r *M20260825000002CreateTopicsTable) Down() error {
	return facades.Schema().DropIfExists("topics")
}
