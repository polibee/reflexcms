package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260824000005CreateTagsTables struct{}

func (r *M20260824000005CreateTagsTables) Signature() string {
	return "20260824000005_create_tags_tables"
}

func (r *M20260824000005CreateTagsTables) Up() error {
	if !facades.Schema().HasTable("tags") {
		if err := facades.Schema().Create("tags", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 128)
			table.String("slug", 191)
			table.TimestampsTz()
			table.Unique("slug")
			table.Unique("name")
		}); err != nil {
			return err
		}
	}

	if facades.Schema().HasTable("article_tags") {
		return nil
	}
	return facades.Schema().Create("article_tags", func(table schema.Blueprint) {
		table.UnsignedBigInteger("article_id")
		table.UnsignedBigInteger("tag_id")
		// Composite uniqueness doubles as the join index for both directions.
		table.Unique("article_id", "tag_id")
		table.Index("tag_id")
	})
}

func (r *M20260824000005CreateTagsTables) Down() error {
	if err := facades.Schema().DropIfExists("article_tags"); err != nil {
		return err
	}
	return facades.Schema().DropIfExists("tags")
}
