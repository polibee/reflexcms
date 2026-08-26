package seeders

import (
	"reflexcms/backend/app/facades"
	accessmodels "reflexcms/backend/modules/access/models"
)

type RoleSeeder struct{}

func (s *RoleSeeder) Signature() string { return "role-seeder" }

func (s *RoleSeeder) Run() error {
	count, err := facades.Orm().Query().Model(&accessmodels.Role{}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Permission naming: {prefix}.{action}; wildcard "*" and "prefix.*" only —
	// suffix patterns like "*.view" are NOT supported by the matcher.
	roles := []accessmodels.Role{
		{
			Name:        "super-admin",
			DisplayName: "超级管理员",
			Permissions: accessmodels.StringList{"*"},
			Sort:        10,
		},
		{
			Name:        "editor",
			DisplayName: "编辑",
			Permissions: accessmodels.StringList{
				"articles.view", "articles.create", "articles.edit", "articles.delete", "articles.publish",
				"categories.view", "categories.create", "categories.edit", "categories.delete",
				"tags.view", "tags.create", "tags.edit", "tags.delete",
				"comments.view",
				"users.view",
			},
			Sort: 20,
		},
		{
			Name:        "moderator",
			DisplayName: "审核员",
			Permissions: accessmodels.StringList{
				"topics.view", "topics.edit", "topics.delete", "topics.pin", "topics.feature", "topics.close",
				"replies.view", "replies.delete",
				"comments.view", "comments.moderate",
				"users.view",
			},
			Sort: 30,
		},
		{
			Name:        "viewer",
			DisplayName: "只读",
			Permissions: accessmodels.StringList{
				"users.view", "roles.view",
				"articles.view", "categories.view", "tags.view",
				"forums.view", "topics.view", "replies.view",
				"comments.view", "notifications.view", "settings.view",
			},
			Sort: 40,
		},
	}

	return facades.Orm().Query().Create(&roles)
}
