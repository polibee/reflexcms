package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260824000004CreateArticleCategoriesTable struct{}

func (r *M20260824000004CreateArticleCategoriesTable) Signature() string {
	return "20260824000004_create_article_categories_table"
}

func (r *M20260824000004CreateArticleCategoriesTable) Up() error {
	if facades.Schema().HasTable("article_categories") {
		return nil
	}
	return facades.Schema().Create("article_categories", func(table schema.Blueprint) {
		table.ID()
		table.String("name", 128)
		table.String("slug", 191)
		table.Text("description").Nullable()
		table.Integer("sort").Default(100)
		table.TimestampsTz()
		table.Unique("slug")
	})
}

func (r *M20260824000004CreateArticleCategoriesTable) Down() error {
	return facades.Schema().DropIfExists("article_categories")
}
