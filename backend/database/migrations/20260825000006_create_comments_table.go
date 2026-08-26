package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000006CreateCommentsTable struct{}

func (r *M20260825000006CreateCommentsTable) Signature() string {
	return "20260825000006_create_comments_table"
}

// Up creates the moderated article-comment store. This is deliberately a
// separate domain from forum replies — see docs/模块化架构与产品组合.md §7.
func (r *M20260825000006CreateCommentsTable) Up() error {
	if facades.Schema().HasTable("comments") {
		return nil
	}
	return facades.Schema().Create("comments", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("article_id")
		table.UnsignedBigInteger("user_id").Nullable()
		table.UnsignedBigInteger("parent_id").Nullable()
		table.Text("content")
		table.String("status", 32).Default("pending")
		table.SoftDeletes()
		table.TimestampsTz()
		table.Index("article_id")
		table.Index("status")
	})
}

func (r *M20260825000006CreateCommentsTable) Down() error {
	return facades.Schema().DropIfExists("comments")
}
