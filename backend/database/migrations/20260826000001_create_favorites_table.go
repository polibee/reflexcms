package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260826000001CreateFavoritesTable struct{}

func (r *M20260826000001CreateFavoritesTable) Signature() string {
	return "20260826000001_create_favorites_table"
}

func (r *M20260826000001CreateFavoritesTable) Up() error {
	if facades.Schema().HasTable("favorites") {
		return nil
	}
	return facades.Schema().Create("favorites", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.String("fav_type", 16).Default("topic")
		table.UnsignedBigInteger("fav_id")
		table.DateTimeTz("created_at").Nullable()
		table.Unique("user_id", "fav_type", "fav_id")
		table.Index("user_id")
	})
}

func (r *M20260826000001CreateFavoritesTable) Down() error {
	return facades.Schema().Drop("favorites")
}
