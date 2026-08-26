package http

import (
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	adminhub "reflexcms/backend/modules/adminhub"
	auditservices "reflexcms/backend/modules/audit/services"
	forummodels "reflexcms/backend/modules/forum/models"
)

// BestReplyController marks / clears the accepted answer of a topic
// (Flarum/Discourse style). POST /api/admin/topics/{id}/best-reply with
// {reply_id}; reply_id=0 clears the marker.
type BestReplyController struct{}

func NewBestReplyController() *BestReplyController { return &BestReplyController{} }

func (r *BestReplyController) SetBestReply(ctx http.Context) http.Response {
	if resp := adminhub.CheckPermission(ctx, "topics.edit"); resp != nil {
		return *resp
	}
	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "topic not found")
	}

	var req struct {
		ReplyID uint64 `json:"reply_id"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "Invalid body")
	}

	var topic forummodels.Topic
	if err := facades.Orm().Query().FindOrFail(&topic, id); err != nil {
		return httpx.Error(ctx, 404, "topic not found")
	}

	if req.ReplyID == 0 {
		topic.BestReplyID = 0
	} else {
		var reply forummodels.Reply
		// The reply must belong to this topic — cross-topic answers are rejected.
		if err := facades.Orm().Query().
			Where("id = ? AND topic_id = ?", req.ReplyID, id).
			FirstOrFail(&reply); err != nil {
			return httpx.Error(ctx, 422, "reply does not belong to this topic")
		}
		topic.BestReplyID = req.ReplyID
	}

	if err := facades.Orm().Query().Save(&topic); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	auditservices.Record(ctx, adminhub.CurrentUserID(ctx), "best-reply", "topics", topic.ID, http.StatusOK,
		map[string]any{"reply_id": req.ReplyID})

	return ctx.Response().Success().Json(http.Json{"id": topic.ID, "best_reply_id": topic.BestReplyID})
}
