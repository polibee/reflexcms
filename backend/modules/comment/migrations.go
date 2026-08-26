package comment

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/database/migrations"
)

func migrationsComment() []schema.Migration {
	return migrations.CommentMigrations()
}
