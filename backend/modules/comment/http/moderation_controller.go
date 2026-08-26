package http

import (
	"github.com/goravel/framework/contracts/http"

	httpx "reflexcms/backend/app/support/httpx"
	adminhub "reflexcms/backend/modules/adminhub"
	auditservices "reflexcms/backend/modules/audit/services"
	commentmodels "reflexcms/backend/modules/comment/models"
	commentservices "reflexcms/backend/modules/comment/services"
)

// ModerationController serves POST /api/admin/comments/{id}/approve|reject.
type ModerationController struct{}

func NewModerationController() *ModerationController { return &ModerationController{} }

func (r *ModerationController) Approve(ctx http.Context) http.Response {
	return r.moderate(ctx, "approve", commentservices.Approve)
}

func (r *ModerationController) Reject(ctx http.Context) http.Response {
	return r.moderate(ctx, "reject", commentservices.Reject)
}

func (r *ModerationController) moderate(ctx http.Context, action string,
	fn func(uint64) (*commentmodels.Comment, error)) http.Response {
	if resp := adminhub.CheckPermission(ctx, "comments.moderate"); resp != nil {
		return *resp
	}

	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "comment not found")
	}

	comment, err := fn(id)
	if err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}

	auditservices.Record(ctx, adminhub.CurrentUserID(ctx), action, "comments", comment.ID, http.StatusOK, nil)

	return ctx.Response().Success().Json(http.Json{
		"id": comment.ID, "status": comment.Status,
	})
}
