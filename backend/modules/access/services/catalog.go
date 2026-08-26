package services

// PermissionCatalog is the single source of truth for grantable permissions.
// The frontend role editor renders checkboxes from the same catalog exposed
// via GET /api/admin/permissions; writes are validated against it so typos
// can never silently grant nothing.

type Permission struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type PermissionGroup struct {
	Group       string       `json:"group"`
	Permissions []Permission `json:"permissions"`
}

var catalog = []PermissionGroup{
	{
		Group: "内容",
		Permissions: []Permission{
			{Name: "articles.view", Label: "文章查看"}, {Name: "articles.create", Label: "文章创建"},
			{Name: "articles.edit", Label: "文章编辑"}, {Name: "articles.delete", Label: "文章删除"},
			{Name: "articles.publish", Label: "文章发布/定时"},
			{Name: "categories.view", Label: "分类查看"}, {Name: "categories.create", Label: "分类创建"},
			{Name: "categories.edit", Label: "分类编辑"}, {Name: "categories.delete", Label: "分类删除"},
			{Name: "tags.view", Label: "标签查看"}, {Name: "tags.create", Label: "标签创建"},
			{Name: "tags.edit", Label: "标签编辑"}, {Name: "tags.delete", Label: "标签删除"},
			{Name: "comments.view", Label: "评论查看"}, {Name: "comments.moderate", Label: "评论审核"},
		},
	},
	{
		Group: "论坛",
		Permissions: []Permission{
			{Name: "forums.view", Label: "板块查看"}, {Name: "forums.create", Label: "板块创建"},
			{Name: "forums.edit", Label: "板块编辑"}, {Name: "forums.delete", Label: "板块删除"},
			{Name: "topics.view", Label: "帖子查看"}, {Name: "topics.edit", Label: "帖子编辑"},
			{Name: "topics.delete", Label: "帖子删除"}, {Name: "topics.pin", Label: "帖子置顶"},
			{Name: "topics.feature", Label: "帖子加精"},
			{Name: "replies.view", Label: "回复查看"}, {Name: "replies.create", Label: "回复创建"},
			{Name: "replies.delete", Label: "回复删除"},
		},
	},
	{
		Group: "平台",
		Permissions: []Permission{
			{Name: "users.view", Label: "用户查看"}, {Name: "users.create", Label: "用户创建"},
			{Name: "users.edit", Label: "用户编辑"}, {Name: "users.delete", Label: "用户删除"},
			{Name: "roles.view", Label: "角色查看"}, {Name: "roles.create", Label: "角色创建"},
			{Name: "roles.edit", Label: "角色编辑"}, {Name: "roles.delete", Label: "角色删除"},
			{Name: "notifications.view", Label: "通知查看"}, {Name: "notifications.edit", Label: "通知操作"},
			{Name: "settings.view", Label: "设置查看"}, {Name: "settings.edit", Label: "设置修改"},
			{Name: "audit.view", Label: "审计日志查看"},
		},
	},
}

// Catalog returns the full grouped permission catalog.
func Catalog() []PermissionGroup {
	return catalog
}

// IsKnownPermission reports whether name is in the catalog. Wildcards
// ("*" and "prefix.*") are always considered valid.
func IsKnownPermission(name string) bool {
	if name == "*" || (len(name) > 1 && name[len(name)-1] == '*') {
		return true
	}
	for _, g := range catalog {
		for _, p := range g.Permissions {
			if p.Name == name {
				return true
			}
		}
	}
	return false
}
