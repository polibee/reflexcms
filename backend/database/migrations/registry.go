package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
)

// Ownership groups (docs/模块化架构与产品组合.md §3): each feature module
// declares which migrations belong to it, so products only ever run the
// schema they need. Physical files stay in this package for uniform naming.

// PlatformMigrations are required by the core regardless of composition.
func PlatformMigrations() []schema.Migration {
	return []schema.Migration{
		&M20210101000001CreateJobsTable{},
	}
}

// AccessMigrations: roles + users (RBAC domain) + per-user favorites.
func AccessMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260824000001CreateRolesTable{},
		&M20260824000002CreateUsersTable{},
		&M20260826000001CreateFavoritesTable{},
		&M20260826000002AddProfileCardToUsers{},
		&M20260826000003CreateMessagesTable{},
		&M20260826000005CreateBlockedUsersTable{},
	}
}

// PagesMigrations: admin-authored static pages (隐私政策、TOS 等).
func PagesMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260826000004CreatePagesTable{},
	}
}

// PayMigrations: products + orders (商城).
func PayMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260826000006CreateProductsTable{},
		&M20260826000007CreateOrdersTable{},
	}
}

// AuthMigrations: admin session storage.
func AuthMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260824000003CreateAdminSessionsTable{},
	}
}

// CmsMigrations: article categories, tags and articles.
func CmsMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260824000004CreateArticleCategoriesTable{},
		&M20260824000005CreateTagsTables{},
		&M20260824000006CreateArticlesTable{},
		&M20260825000009EnhanceCms{},
	}
}

// ForumMigrations: boards, topics, replies, likes.
func ForumMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000001CreateForumCategoriesTable{},
		&M20260825000002CreateTopicsTable{},
		&M20260825000003CreateRepliesTable{},
		&M20260825000004CreateLikesTable{},
		&M20260825000010EnhanceForum{},
		&M20260826000008CreateBoardModeratorsTable{},
	}
}

// SettingsMigrations: site KV configuration.
func SettingsMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000005CreateSettingsTable{},
	}
}

// CommentMigrations: moderated article comments.
func CommentMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000006CreateCommentsTable{},
	}
}

// NotificationMigrations: in-app notifications.
func NotificationMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000007CreateNotificationsTable{},
	}
}

// AuditMigrations: admin operation trail.
func AuditMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000008CreateAdminOperationLogsTable{},
	}
}

// LayoutMigrations: widgets, menus, carousels, invite codes, points.
func LayoutMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000011CreateWidgetsTable{},
		&M20260825000012CreateMenusTable{},
		&M20260825000013CreateCarouselsTable{},
		&M20260825000014CreateInviteCodesTable{},
		&M20260825000015CreatePointsTables{},
	}
}

// PointsMigrations: user points balance + transaction history.
func PointsMigrations() []schema.Migration {
	return []schema.Migration{
		&M20260825000015CreatePointsTables{},
	}
}
