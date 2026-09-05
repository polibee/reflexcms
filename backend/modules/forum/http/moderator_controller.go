package http

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	accessmatcher "reflexcms/backend/modules/access/services"
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

	// Governance linkage: super-admin always passes; users carrying the
	// moderator ROLE (RBAC: topics.pin etc.) moderate site-wide; everyone
	// else needs a per-board assignment in board_moderators.
	canModerate := identity.Role == "super-admin"
	if !canModerate && accessmatcher.Can(identity.Permissions, "topics.pin") {
		canModerate = true
	}
	if !canModerate {
		canModerate = IsModerator(boardID, identity.ID)
	}
	if !canModerate {
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

// callerModeratesAny reports whether the identity moderates at least one
// board (site-wide governance floor for mute/ban).
func callerModeratesAny(userID uint64) bool {
	var count int64
	count, _ = facades.Orm().Query().Table("board_moderators").
		Where("user_id = ?", userID).Count()
	return count > 0
}

// MuteUser: POST /api/v1/users/{id}/mute — {duration_hours}. 1–875999
// hours (~100 years = permanent). Moderators only; super-admin bypasses.
func MuteUser(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}
	if identity.Role != "super-admin" && !callerModeratesAny(identity.ID) {
		return httpx.Error(ctx, 403, "仅版主可以执行禁言")
	}

	targetID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil || targetID == 0 {
		return httpx.Error(ctx, 404, "user not found")
	}
	if targetID == identity.ID {
		return httpx.Error(ctx, 422, "不能禁言自己")
	}

	var req struct {
		DurationHours int `json:"duration_hours"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.DurationHours < 1 || req.DurationHours > 875999 {
		return httpx.Error(ctx, 422, "禁言时长无效（1 小时 ~ 永久）")
	}

	hours := req.DurationHours
	message := "已禁言"
	if hours >= 875999 {
		// permanent: far-future timestamp
		if _, err := facades.Orm().Query().Exec(`
			UPDATE users SET muted_until = '2999-01-01T00:00:00Z'
			WHERE id = ? AND deleted_at IS NULL
		`, targetID); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		message = "已永久禁言"
	} else {
		if _, err := facades.Orm().Query().Exec(`
			UPDATE users SET muted_until = NOW() + make_interval(hours => ?)
			WHERE id = ? AND deleted_at IS NULL
		`, hours, targetID); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		message = fmt.Sprintf("已禁言 %d 小时", hours)
	}
	return ctx.Response().Success().Json(http.Json{"message": message})
}

// UnmuteUser: POST /api/v1/users/{id}/unmute — lifts the mute early.
func UnmuteUser(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}
	if identity.Role != "super-admin" && !callerModeratesAny(identity.ID) {
		return httpx.Error(ctx, 403, "仅版主可以解除禁言")
	}
	targetID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil || targetID == 0 {
		return httpx.Error(ctx, 404, "user not found")
	}
	if _, err := facades.Orm().Query().Exec(`
		UPDATE users SET muted_until = NULL WHERE id = ?
	`, targetID); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"message": "已解除禁言"})
}

// BanUser: POST /api/v1/users/{id}/ban — {days}. Sets banned_until; login
// auto-restores the account once the term passes. Super-admin only — bans
// lock accounts out entirely, so they are not delegated to moderators.
func BanUser(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}
	if identity.Role != "super-admin" {
		return httpx.Error(ctx, 403, "仅超级管理员可以封禁账号")
	}
	targetID, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil || targetID == 0 {
		return httpx.Error(ctx, 404, "user not found")
	}
	if targetID == identity.ID {
		return httpx.Error(ctx, 422, "不能封禁自己")
	}

	var req struct {
		Days int `json:"days"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Days < 1 || req.Days > 365 {
		return httpx.Error(ctx, 422, "封禁天数须在 1–365 天之间")
	}

	if _, err := facades.Orm().Query().Exec(`
		UPDATE users
		SET banned_until = NOW() + make_interval(days => ?), status = 'banned', banned_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`, req.Days, targetID); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{
		"message": fmt.Sprintf("已封禁 %d 天", req.Days),
	})
}
