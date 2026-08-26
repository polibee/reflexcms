package http

import (
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	adminhub "reflexcms/backend/modules/adminhub"
	auditservices "reflexcms/backend/modules/audit/services"
	forummodels "reflexcms/backend/modules/forum/models"
)

// TopicActionController serves the moderation routes:
//
//	POST /api/admin/topics/{id}/pin|unpin|feature|unfeature|close|open
type TopicActionController struct{}

func NewTopicActionController() *TopicActionController {
	return &TopicActionController{}
}

func (r *TopicActionController) Pin(ctx http.Context) http.Response {
	return r.setFlag(ctx, "pin", true, "topics.pin")
}

func (r *TopicActionController) Unpin(ctx http.Context) http.Response {
	return r.setFlag(ctx, "pin", false, "topics.pin")
}

func (r *TopicActionController) Feature(ctx http.Context) http.Response {
	return r.setFlag(ctx, "feature", true, "topics.feature")
}

func (r *TopicActionController) Unfeature(ctx http.Context) http.Response {
	return r.setFlag(ctx, "feature", false, "topics.feature")
}

func (r *TopicActionController) Close(ctx http.Context) http.Response {
	return r.setStatus(ctx, forummodels.TopicClosed, "topics.edit")
}

func (r *TopicActionController) Open(ctx http.Context) http.Response {
	return r.setStatus(ctx, forummodels.TopicOpen, "topics.edit")
}

func (r *TopicActionController) setFlag(ctx http.Context, name string, value bool, permission string) http.Response {
	if resp := adminhub.CheckPermission(ctx, permission); resp != nil {
		return *resp
	}
	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "topic not found")
	}

	var topic forummodels.Topic
	if err := facades.Orm().Query().FindOrFail(&topic, id); err != nil {
		return httpx.Error(ctx, 404, "topic not found")
	}

	switch name {
	case "pin":
		topic.IsPinned = value
	case "feature":
		topic.IsFeatured = value
	}
	if err := facades.Orm().Query().Save(&topic); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	auditservices.Record(ctx, adminhub.CurrentUserID(ctx), name, "topics", topic.ID, http.StatusOK, nil)

	return ctx.Response().Success().Json(http.Json{
		"id": topic.ID, "is_pinned": topic.IsPinned, "is_featured": topic.IsFeatured,
	})
}

func (r *TopicActionController) setStatus(ctx http.Context, status string, permission string) http.Response {
	if resp := adminhub.CheckPermission(ctx, permission); resp != nil {
		return *resp
	}
	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "topic not found")
	}

	var topic forummodels.Topic
	if err := facades.Orm().Query().FindOrFail(&topic, id); err != nil {
		return httpx.Error(ctx, 404, "topic not found")
	}

	topic.Status = status
	if err := facades.Orm().Query().Save(&topic); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	auditservices.Record(ctx, adminhub.CurrentUserID(ctx), status, "topics", topic.ID, http.StatusOK, nil)

	return ctx.Response().Success().Json(http.Json{"id": topic.ID, "status": topic.Status})
}
