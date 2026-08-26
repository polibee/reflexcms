package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000001CreateForumCategoriesTable struct{}

func (r *M20260825000001CreateForumCategoriesTable) Signature() string {
	return "20260825000001_create_forum_categories_table"
}

func (r *M20260825000001CreateForumCategoriesTable) Up() error {
	if facades.Schema().HasTable("forum_categories") {
		return nil
	}
	return facades.Schema().Create("forum_categories", func(table schema.Blueprint) {
		table.ID()
		table.String("name", 128)
		table.String("slug", 191)
		table.Text("description").Nullable()
		table.String("icon", 64).Nullable()
		table.Integer("sort").Default(100)
		// Denormalised topic counter, maintained in transactions.
		table.UnsignedBigInteger("topic_count").Default(0)
		table.TimestampsTz()
		table.Unique("slug")
	})
}

func (r *M20260825000001CreateForumCategoriesTable) Down() error {
	return facades.Schema().DropIfExists("forum_categories")
}
