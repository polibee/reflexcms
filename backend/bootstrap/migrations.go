package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260824000001CreateRolesTable{},
		&migrations.M20260824000002CreateUsersTable{},
		&migrations.M20260824000003CreateAdminSessionsTable{},
		&migrations.M20260824000004CreateArticleCategoriesTable{},
		&migrations.M20260824000005CreateTagsTables{},
		&migrations.M20260824000006CreateArticlesTable{},
		&migrations.M20260825000001CreateForumCategoriesTable{},
		&migrations.M20260825000002CreateTopicsTable{},
		&migrations.M20260825000003CreateRepliesTable{},
		&migrations.M20260825000004CreateLikesTable{},
	}
}
