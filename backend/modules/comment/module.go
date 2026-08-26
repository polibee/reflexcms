package comment

import (
	"github.com/goravel/framework/contracts/database/schema"

	adminhub "reflexcms/backend/modules/adminhub"
	commenthttp "reflexcms/backend/modules/comment/http"
	commentmodels "reflexcms/backend/modules/comment/models"
)

// Module is the article-comment feature module: pending-by-default comments
// with a moderation state machine. Distinct from forum replies by design
// (docs/模块化架构与产品组合.md §7).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "comment" }

func (m *Module) Migrations() []schema.Migration {
	return migrationsComment()
}

func (m *Module) Routes() {
	adminhub.Register(&adminhub.Spec{
		Name:       "comments",
		Model:      commentmodels.Comment{},
		Searchable: []string{"content"},
		Sortable:   []string{"id", "status", "article_id", "created_at"},
		// status is NOT fillable: it changes only through moderation actions.
		Fillable:         []string{"article_id", "user_id", "parent_id", "content"},
		PermissionPrefix: "comments",

		SoftDeletable: true,
	})

	moderation := commenthttp.NewModerationController()
	publicAPI := commenthttp.NewPublicCommentController()
	r := facadesRoute()

	r.Post("/api/admin/comments/{id}/approve", moderation.Approve)
	r.Post("/api/admin/comments/{id}/reject", moderation.Reject)
	r.Get("/api/v1/articles/{slug}/comments", publicAPI.List)
	r.Post("/api/v1/articles/{slug}/comments", publicAPI.Submit)

	// The {id}/action subtrees shadow the wildcard bulk-delete route; re-bind
	// the static path for this resource (see adminhub.RegisterBulkDeleteFor).
	adminhub.RegisterBulkDeleteFor("comments")
}

func (m *Module) Boot() {}
