package http

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	authhttp "reflexcms/backend/modules/auth/http/controllers"
)

/* Board-moderator powers: v1 moderation actions gated by per-board
 * assignment (board_moderators), not by admin RBAC. Moderators can pin /
 * feature / close / open / best-reply inside their own boards. */

// IsModerator reports whether userID moderates the given board.
func IsModerator(boardID, userID uint64) bool {
	if boardID == 0 || userID == 0 {
		return false
	}
	var count int64
	count, _ = facades.Orm().Query().Table("board_moderators").
		Where("board_id = ? AND user_id = ?", boardID, userID).Count()
	return count > 0
}

// ModerateTopic: POST /api/v1/topics/{id}/moderate — body {action,
// reply_id?}. Allowed actions: pin/unpin/feature/unfeature/close/open/
// best_reply.
func (r *PublicController) ModerateTopic(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil || id == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}

	var req struct {
		Action  string `json:"action"`
		ReplyID uint64 `json:"reply_id"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Action == "" {
		return httpx.Error(ctx, 422, "action is required")
	}

	var topics []map[string]any
	if err := facades.Orm().Query().Table("topics").
		Where("id = ? AND deleted_at IS NULL", id).
		Select("id", "forum_category_id", "best_reply_id").
		Get(&topics); err != nil || len(topics) == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}
	boardID, _ := toUint64(topics[0]["forum_category_id"])

	// Admins (forum.manage) bypass the board assignment; moderators need it.
	if identity.Role != "super-admin" && !IsModerator(boardID, identity.ID) {
		return httpx.Error(ctx, 403, "只有该板块的版主才能执行此操作")
	}

	updates := map[string]any{}
	switch req.Action {
	case "pin":
		updates["is_pinned"] = true
	case "unpin":
		updates["is_pinned"] = false
	case "feature":
		updates["is_featured"] = true
	case "unfeature":
		updates["is_featured"] = false
	case "close":
		updates["status"] = "closed"
	case "open":
		updates["status"] = "open"
	case "best_reply":
		if req.ReplyID == 0 {
			return httpx.Error(ctx, 422, "reply_id is required for best_reply")
		}
		updates["best_reply_id"] = req.ReplyID
	default:
		return httpx.Error(ctx, 422, fmt.Sprintf("unknown action %q", req.Action))
	}

	if _, err := facades.Orm().Query().Table("topics").
		Where("id = ?", id).Update(updates); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"message": "操作成功", "action": req.Action})
}
