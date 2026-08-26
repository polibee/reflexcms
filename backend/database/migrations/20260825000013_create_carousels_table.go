package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000013CreateCarouselsTable struct{}

func (r *M20260825000013CreateCarouselsTable) Signature() string {
	return "20260825000013_create_carousels_table"
}

func (r *M20260825000013CreateCarouselsTable) Up() error {
	if facades.Schema().HasTable("carousels") {
		return nil
	}
	return facades.Schema().Create("carousels", func(table schema.Blueprint) {
		table.ID()
		table.String("title", 255).Nullable()
		table.String("image_url", 1024)
		table.String("link_url", 1024).Nullable()
		table.Integer("sort").Default(100)
		table.Boolean("is_active").Default(true)
		table.TimestampsTz()
		table.Index("is_active")
	})
}

func (r *M20260825000013CreateCarouselsTable) Down() error {
	return facades.Schema().DropIfExists("carousels")
}
