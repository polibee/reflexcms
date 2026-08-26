package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260826000004CreatePagesTable struct{}

func (r *M20260826000004CreatePagesTable) Signature() string {
	return "20260826000004_create_pages_table"
}

func (r *M20260826000004CreatePagesTable) Up() error {
	if facades.Schema().HasTable("pages") {
		return nil
	}
	return facades.Schema().Create("pages", func(table schema.Blueprint) {
		table.ID()
		table.String("title", 255)
		table.String("slug", 191)
		table.LongText("content").Nullable()
		table.String("status", 32).Default("draft")
		table.TimestampsTz()
		table.Unique("slug")
		table.Index("status")
	})
}

func (r *M20260826000004CreatePagesTable) Down() error {
	return facades.Schema().Drop("pages")
}
