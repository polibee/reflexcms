package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000010EnhanceForum struct{}

func (r *M20260825000010EnhanceForum) Signature() string {
	return "20260825000010_enhance_forum"
}

// Up adds the Flarum/Discourse-inspired "best reply" marker.
func (r *M20260825000010EnhanceForum) Up() error {
	if facades.Schema().HasColumn("topics", "best_reply_id") {
		return nil
	}
	return facades.Schema().Table("topics", func(table schema.Blueprint) {
		table.UnsignedBigInteger("best_reply_id").Nullable()
	})
}

func (r *M20260825000010EnhanceForum) Down() error {
	if !facades.Schema().HasColumn("topics", "best_reply_id") {
		return nil
	}
	return facades.Schema().Table("topics", func(table schema.Blueprint) {
		table.DropColumn("best_reply_id")
	})
}
