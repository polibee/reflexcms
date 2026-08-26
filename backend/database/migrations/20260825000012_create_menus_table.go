package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000012CreateMenusTable struct{}

func (r *M20260825000012CreateMenusTable) Signature() string {
	return "20260825000012_create_menus_table"
}

func (r *M20260825000012CreateMenusTable) Up() error {
	if facades.Schema().HasTable("menus") {
		return nil
	}
	return facades.Schema().Create("menus", func(table schema.Blueprint) {
		table.ID()
		table.String("location", 32)
		table.String("label", 128)
		table.String("url", 512)
		table.Integer("sort").Default(100)
		table.UnsignedBigInteger("parent_id").Nullable()
		table.Boolean("is_visible").Default(true)
		table.TimestampsTz()
		table.Index("location")
	})
}

func (r *M20260825000012CreateMenusTable) Down() error {
	return facades.Schema().DropIfExists("menus")
}
