package migrations

import (
	"fmt"

	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260824000006CreateArticlesTable struct{}

func (r *M20260824000006CreateArticlesTable) Signature() string {
	return "20260824000006_create_articles_table"
}

func (r *M20260824000006CreateArticlesTable) Up() error {
	if !facades.Schema().HasTable("articles") {
		if err := facades.Schema().Create("articles", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("author_id").Nullable()
			table.UnsignedBigInteger("category_id").Nullable()
			table.String("title", 512)
			table.String("slug", 191)
			// Markdown source; rendered HTML belongs to the (future) frontend.
			table.LongText("content").Nullable()
			table.Text("summary").Nullable()
			table.String("cover", 512).Nullable()
			table.String("status", 32).Default("draft")
			table.DateTimeTz("published_at").Nullable()
			table.UnsignedBigInteger("view_count").Default(0)
			table.UnsignedBigInteger("comment_count").Default(0)
			table.Text("ai_summary").Nullable()
			// ai_keywords: JSONB string array, reserved for the AI module.
			table.Jsonb("ai_keywords")
			table.SoftDeletes()
			table.TimestampsTz()
			table.Unique("slug")
			table.Index("status")
			table.Index("author_id")
			table.Index("category_id")
			table.Index("published_at")
		}); err != nil {
			return fmt.Errorf("schema create articles: %w", err)
		}
	}

	// NOTE: raw DDL goes through the ORM connection on purpose —
	// facades.DB() resolves a different database handle than facades.Orm()
	// in this setup, which silently targeted the wrong database.
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE articles DROP COLUMN IF EXISTS search_vector
	`); err != nil {
		return err
	}
	// Generated full-text column: title weighted A, summary weighted B.
	// 'simple' config keeps CJK text as-is; CJK queries degrade to ILIKE in
	// the search layer (plan §M2). No external input is involved here — pure
	// DDL, so raw SQL is safe and binding rules do not apply.
	if _, err := facades.Orm().Query().Exec(`
		ALTER TABLE articles ADD COLUMN search_vector tsvector
		GENERATED ALWAYS AS (
			setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
			setweight(to_tsvector('simple', coalesce(summary, '')), 'B')
		) STORED
	`); err != nil {
		return err
	}
	_, err := facades.Orm().Query().Exec(`
		CREATE INDEX IF NOT EXISTS idx_articles_search_vector
		ON articles USING GIN (search_vector)
	`)
	return err
}

func (r *M20260824000006CreateArticlesTable) Down() error {
	return facades.Schema().DropIfExists("articles")
}
