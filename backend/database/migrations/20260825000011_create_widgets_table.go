package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000011CreateWidgetsTable struct{}

func (r *M20260825000011CreateWidgetsTable) Signature() string {
	return "20260825000011_create_widgets_table"
}

func (r *M20260825000011CreateWidgetsTable) Up() error {
	if facades.Schema().HasTable("widgets") {
		return nil
	}
	return facades.Schema().Create("widgets", func(table schema.Blueprint) {
		table.ID()
		table.String("area", 32)
		table.String("widget_type", 64)
		table.String("title", 255).Nullable()
		table.Jsonb("config")
		table.Integer("sort").Default(100)
		table.Boolean("is_active").Default(true)
		table.TimestampsTz()
		table.Index("area")
	})
}

func (r *M20260825000011CreateWidgetsTable) Down() error {
	return facades.Schema().DropIfExists("widgets")
}
