package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000009EnhanceCms struct{}

func (r *M20260825000009EnhanceCms) Signature() string {
	return "20260825000009_enhance_cms"
}

// Up adds WordPress-inspired capabilities: sticky articles and hierarchical
// categories.
func (r *M20260825000009EnhanceCms) Up() error {
	if !facades.Schema().HasColumn("articles", "is_pinned") {
		if err := facades.Schema().Table("articles", func(table schema.Blueprint) {
			table.Boolean("is_pinned").Default(false)
			table.Index("is_pinned")
		}); err != nil {
			return err
		}
	}

	if !facades.Schema().HasColumn("article_categories", "parent_id") {
		if err := facades.Schema().Table("article_categories", func(table schema.Blueprint) {
			table.UnsignedBigInteger("parent_id").Nullable()
			table.Index("parent_id")
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260825000009EnhanceCms) Down() error {
	if facades.Schema().HasColumn("articles", "is_pinned") {
		if err := facades.Schema().Table("articles", func(table schema.Blueprint) {
			table.DropColumn("is_pinned")
		}); err != nil {
			return err
		}
	}
	if facades.Schema().HasColumn("article_categories", "parent_id") {
		if err := facades.Schema().Table("article_categories", func(table schema.Blueprint) {
			table.DropColumn("parent_id")
		}); err != nil {
			return err
		}
	}
	return nil
}
