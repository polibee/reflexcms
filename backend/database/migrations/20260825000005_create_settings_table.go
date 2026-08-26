package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000005CreateSettingsTable struct{}

func (r *M20260825000005CreateSettingsTable) Signature() string {
	return "20260825000005_create_settings_table"
}

func (r *M20260825000005CreateSettingsTable) Up() error {
	if facades.Schema().HasTable("settings") {
		return nil
	}
	return facades.Schema().Create("settings", func(table schema.Blueprint) {
		table.ID()
		table.String("key", 191)
		// JSON-encoded scalar/object; the service layer owns the codec.
		table.Text("value").Nullable()
		table.String("group", 64).Default("general")
		table.TimestampsTz()
		table.Unique("key")
	})
}

func (r *M20260825000005CreateSettingsTable) Down() error {
	return facades.Schema().DropIfExists("settings")
}
