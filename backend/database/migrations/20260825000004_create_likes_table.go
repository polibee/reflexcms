package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000004CreateLikesTable struct{}

func (r *M20260825000004CreateLikesTable) Signature() string {
	return "20260825000004_create_likes_table"
}

func (r *M20260825000004CreateLikesTable) Up() error {
	if facades.Schema().HasTable("likes") {
		return nil
	}
	return facades.Schema().Create("likes", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.String("likeable_type", 64)
		table.UnsignedBigInteger("likeable_id")
		table.TimestampsTz()
		// The composite unique index makes like toggles idempotent.
		table.Unique("user_id", "likeable_type", "likeable_id")
		table.Index("likeable_type", "likeable_id")
	})
}

func (r *M20260825000004CreateLikesTable) Down() error {
	return facades.Schema().DropIfExists("likes")
}
