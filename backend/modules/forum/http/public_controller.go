package http

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	"reflexcms/backend/modules/access"
	authhttp "reflexcms/backend/modules/auth/http/controllers"
	forumservices "reflexcms/backend/modules/forum/services"
	notificationservices "reflexcms/backend/modules/notification/services"
	settingsservices "reflexcms/backend/modules/settings/services"
)

// PublicController serves anonymous /api/v1 read endpoints for the Forum
// domain, gated by site.mode.
type PublicController struct{}

func NewPublicController() *PublicController { return &PublicController{} }

func gateOr404(ctx http.Context, section string) *http.Response {
	if !settingsservices.SectionEnabled(section) {
		resp := httpx.Error(ctx, 404, "Not Found")
		return &resp
	}
	return nil
}

// ListForums: GET /api/v1/forums — all boards with topic counts.
func (r *PublicController) ListForums(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "forums"); resp != nil {
		return *resp
	}

	items := []map[string]any{}
	if err := facades.Orm().Query().Table("forum_categories").
		Select("id", "name", "slug", "description", "icon", "sort", "topic_count").
		OrderBy("sort", "asc").
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	return ctx.Response().Success().Json(http.Json{"items": items})
}

// ListTopics: GET /api/v1/topics?forum_category_id=&sort=latest|hot&page&perPage
func (r *PublicController) ListTopics(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}

	page := clampPageV1(ctx.Request().Query("page", "1"))
	perPage := clampPerPageV1(ctx.Request().Query("perPage", "20"))

	query := facades.Orm().Query().Table("topics").
		Where("status = 'open' AND deleted_at IS NULL").
		Select("id", "title", "forum_category_id", "user_id", "last_reply_user_id",
			"is_pinned", "is_featured", "reply_count", "like_count", "view_count",
			"last_reply_at", "created_at")

	if fcID := ctx.Request().Query("forum_category_id", ""); fcID != "" {
		query = query.Where("forum_category_id = ?", fcID)
	}

	switch ctx.Request().Query("sort", "latest") {
	case "hot":
		query = query.OrderBy("(view_count*1 + reply_count*3 + like_count*5)", "desc")
	case "new":
		// Newest topics first, regardless of activity.
		query = query.OrderBy("created_at", "desc")
	case "reply":
		// Most recently replied topics (activity feed).
		query = query.OrderBy("last_reply_at", "desc").OrderBy("id", "desc")
	default:
		query = query.OrderBy("is_pinned", "desc").OrderBy("last_reply_at", "desc")
	}

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.Offset((page - 1) * perPage).Limit(perPage).Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	attachBoardNames(items)
	attachUserNames(items)

	return envelope(ctx, items, int(total), page, perPage)
}

// attachBoardNames fills in board_name for topic rows so public lists can
// show which board each topic belongs to without a SQL join in the builder.
func attachBoardNames(items []map[string]any) {
	ids := map[uint64]bool{}
	for _, it := range items {
		switch v := it["forum_category_id"].(type) {
		case uint64:
			if v > 0 {
				ids[v] = true
			}
		case int64:
			if v > 0 {
				ids[uint64(v)] = true
			}
		}
	}
	if len(ids) == 0 {
		return
	}

	args := make([]any, 0, len(ids))
	placeholders := make([]string, 0, len(ids))
	for id := range ids {
		args = append(args, id)
		placeholders = append(placeholders, "?")
	}

	var boards []map[string]any
	if err := facades.Orm().Query().Table("forum_categories").
		Where("id IN ("+strings.Join(placeholders, ", ")+")", args...).
		Select("id", "name").
		Get(&boards); err != nil {
		return
	}

	names := map[uint64]string{}
	for _, b := range boards {
		if id, ok := toUint64(b["id"]); ok {
			names[id], _ = b["name"].(string)
		}
	}
	for _, it := range items {
		if id, ok := toUint64(it["forum_category_id"]); ok {
			if name, ok := names[id]; ok {
				it["board_name"] = name
			}
		}
	}
}

// attachUserNames resolves topic author and last-replier usernames in one
// batched lookup, attaching author_name / last_replier_name / author_signature.
func attachUserNames(items []map[string]any) {
	ids := map[uint64]bool{}
	for _, it := range items {
		if id, ok := toUint64(it["user_id"]); ok && id > 0 {
			ids[id] = true
		}
		if id, ok := toUint64(it["last_reply_user_id"]); ok && id > 0 {
			ids[id] = true
		}
	}
	if len(ids) == 0 {
		return
	}

	args := make([]any, 0, len(ids))
	placeholders := make([]string, 0, len(ids))
	for id := range ids {
		args = append(args, id)
		placeholders = append(placeholders, "?")
	}

	var users []map[string]any
	if err := facades.Orm().Query().Table("users").
		Where("id IN ("+strings.Join(placeholders, ", ")+")", args...).
		Select("id", "username", "signature").
		Get(&users); err != nil {
		return
	}

	names := map[uint64]string{}
	sigs := map[uint64]string{}
	for _, u := range users {
		if id, ok := toUint64(u["id"]); ok {
			names[id], _ = u["username"].(string)
			sigs[id], _ = u["signature"].(string)
		}
	}
	for _, it := range items {
		if id, ok := toUint64(it["user_id"]); ok {
			if name := names[id]; name != "" {
				it["author_name"] = name
			}
			// attach signatures for reply-style rows (they render under cards)
			if _, isReply := it["content"]; isReply {
				it["author_signature"] = sigs[id]
			}
		}
		if id, ok := toUint64(it["last_reply_user_id"]); ok && names[id] != "" {
			it["last_replier_name"] = names[id]
		}
	}
}

func toUint64(v any) (uint64, bool) {
	switch n := v.(type) {
	case uint64:
		return n, true
	case int64:
		return uint64(n), true
	case int:
		return uint64(n), true
	case uint:
		return uint64(n), true
	default:
		return 0, false
	}
}

// LatestReplies: GET /api/v1/replies/latest — newest replies across all open
// topics, with topic title for the forum activity feed. Paginated.
func (r *PublicController) LatestReplies(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}

	page := clampPageV1(ctx.Request().Query("page", "1"))
	perPage := clampPerPageV1(ctx.Request().Query("perPage", "20"))

	base := facades.Orm().Query().Table("replies as r").
		Join("JOIN topics t ON t.id = r.topic_id AND t.deleted_at IS NULL AND t.status = 'open'").
		Where("r.deleted_at IS NULL")
	// Signed-in callers never see replies from users they blocked.
	if blocked := access.BlockedIDs(ctx); len(blocked) > 0 {
		args := make([]any, 0, len(blocked))
		ph := make([]string, 0, len(blocked))
		for _, b := range blocked {
			args = append(args, b)
			ph = append(ph, "?")
		}
		base = base.Where("r.user_id NOT IN ("+strings.Join(ph, ", ")+")", args...)
	}

	total, err := base.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := base.
		Select("r.id", "r.topic_id", "t.title AS topic_title", "r.user_id", "r.content", "r.floor", "r.created_at").
		OrderBy("r.created_at", "desc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	return envelope(ctx, items, int(total), page, perPage)
}

// ShowTopic: GET /api/v1/topics/:id — topic detail with recent replies.
func (r *PublicController) ShowTopic(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}

	idStr := ctx.Request().Route("id")
	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr != nil || id == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}

	var topic []map[string]any
	if err := facades.Orm().Query().Table("topics").
		Where("id = ? AND status = 'open' AND deleted_at IS NULL", id).
		Limit(1).
		Get(&topic); err != nil || len(topic) == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}

	facades.Cache().Increment("views:topic:" + idStr)

	replies := []map[string]any{}
	replyQuery := facades.Orm().Query().Table("replies").
		Where("topic_id = ? AND deleted_at IS NULL", id)
	// Signed-in callers never see replies from users they blocked.
	if blocked := access.BlockedIDs(ctx); len(blocked) > 0 {
		args := make([]any, 0, len(blocked))
		ph := make([]string, 0, len(blocked))
		for _, b := range blocked {
			args = append(args, b)
			ph = append(ph, "?")
		}
		replyQuery = replyQuery.Where("user_id NOT IN ("+strings.Join(ph, ", ")+")", args...)
	}
	if err := replyQuery.
		Select("id", "user_id", "parent_id", "content", "floor", "like_count", "created_at").
		OrderBy("floor", "asc").Limit(10).
		Get(&replies); err != nil {
		replies = []map[string]any{}
	}

	result := topic[0]
	result["replies"] = replies

	// Frontend moderation affordances for board moderators / super-admin.
	if identity, ok := authhttp.RequireIdentity(ctx); ok {
		canModerate := identity.Role == "super-admin"
		if !canModerate {
			if boardID, ok := toUint64(result["forum_category_id"]); ok && IsModerator(boardID, identity.ID) {
				canModerate = true
			}
		}
		result["can_moderate"] = canModerate
	}

	return ctx.Response().Success().Json(result)
}

// topicThrottleWindow is the minimum interval between two topics from the
// same user; a cheap spam brake on top of session auth.
const topicThrottleWindow = 30 * time.Second

// rejectedByMute checks the session user's mute window; returns a 403
// response when silenced, nil otherwise.
func rejectedByMute(ctx http.Context, userID uint64) *http.Response {
	var until []any
	_ = facades.Orm().Query().Table("users").
		Where("id = ? AND muted_until IS NOT NULL AND muted_until > NOW()", userID).
		Pluck("muted_until", &until)
	if len(until) > 0 {
		resp := httpx.Error(ctx, 403, "你已被禁言，暂时无法发言")
		return &resp
	}
	return nil
}

// StoreTopic: POST /api/v1/topics — creates a new topic for the session
// user. Identity comes from the bearer token; status defaults to open.
func (r *PublicController) StoreTopic(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}

	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录后再发帖")
	}
	if resp := rejectedByMute(ctx, identity.ID); resp != nil {
		return *resp
	}

	var req struct {
		Title           string `json:"title"`
		Content         string `json:"content"`
		ForumCategoryID uint64 `json:"forum_category_id"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}

	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if title == "" {
		return httpx.Error(ctx, 422, "标题不能为空")
	}
	if len([]rune(title)) > 150 {
		return httpx.Error(ctx, 422, "标题过长（最多 150 字）")
	}
	if content == "" {
		return httpx.Error(ctx, 422, "正文不能为空")
	}
	if len(content) > 50000 {
		return httpx.Error(ctx, 422, "正文过长（最多 50000 字）")
	}
	if req.ForumCategoryID == 0 {
		return httpx.Error(ctx, 422, "请选择板块")
	}

	var boardCount int64
	boardCount, _ = facades.Orm().Query().Table("forum_categories").
		Where("id = ?", req.ForumCategoryID).Count()
	if boardCount == 0 {
		return httpx.Error(ctx, 422, "板块不存在")
	}

	throttleKey := "topic:throttle:" + strconv.FormatUint(identity.ID, 10)
	if parseCount(facades.Cache().Get(throttleKey)) > 0 {
		return httpx.Error(ctx, 429, "发帖太快了，请稍等片刻再试")
	}
	_ = facades.Cache().Put(throttleKey, 1, topicThrottleWindow)

	if _, err := facades.Orm().Query().Exec(`
		INSERT INTO topics (user_id, forum_category_id, title, content, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'open', NOW(), NOW())
	`, identity.ID, req.ForumCategoryID, title, content); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	// Postgres Exec can't surface RETURNING ids through the orm facade; the
	// per-user throttle guarantees no same-user race, so the newest topic of
	// this user is the one just inserted.
	var ids []int64
	if err := facades.Orm().Query().Table("topics").
		Where("user_id = ?", identity.ID).
		OrderBy("id", "desc").Limit(1).
		Pluck("id", &ids); err != nil || len(ids) == 0 {
		return httpx.Error(ctx, 500, "topic created but id lookup failed")
	}

	// Daily-capped earn reward (best-effort; never blocks the post).
	_, _ = settingsservices.AwardDaily(identity.ID, settingsservices.PointReasonTopic,
		settingsservices.EarnTopic(), settingsservices.TopicDailyCap())

	return ctx.Response().Success().Json(http.Json{
		"id":      ids[0],
		"message": "发布成功",
	})
}

// TopicReplies: GET /api/v1/topics/{id}/replies — paginated reply list for
// the topic detail page (page 1 is embedded in ShowTopic; deeper pages load
// through here).
func (r *PublicController) TopicReplies(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}

	idStr := ctx.Request().Route("id")
	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr != nil || id == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}

	page := clampPageV1(ctx.Request().Query("page", "1"))
	perPage := clampPerPageV1(ctx.Request().Query("perPage", "20"))

	query := facades.Orm().Query().Table("replies").
		Where("topic_id = ? AND deleted_at IS NULL", id)
	// Signed-in callers never see replies from users they blocked.
	if blocked := access.BlockedIDs(ctx); len(blocked) > 0 {
		args := make([]any, 0, len(blocked))
		ph := make([]string, 0, len(blocked))
		for _, b := range blocked {
			args = append(args, b)
			ph = append(ph, "?")
		}
		query = query.Where("user_id NOT IN ("+strings.Join(ph, ", ")+")", args...)
	}

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.
		Select("id", "user_id", "parent_id", "content", "floor", "like_count", "created_at").
		OrderBy("floor", "asc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return envelope(ctx, items, int(total), page, perPage)
}

// replyThrottleWindow is the minimum interval between two replies from the
// same user; a cheap spam brake on top of session auth.
const replyThrottleWindow = 5 * time.Second

// StoreReply: POST /api/v1/topics/{id}/replies — creates a reply attributed
// to the session user. Identity always comes from the bearer token; anything
// the client claims about user_id is ignored. Goes through ReplyToTopic so
// floors and denormalised counters stay transactionally consistent.
func (r *PublicController) StoreReply(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "topics"); resp != nil {
		return *resp
	}

	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录后再回复")
	}
	if resp := rejectedByMute(ctx, identity.ID); resp != nil {
		return *resp
	}

	idStr := ctx.Request().Route("id")
	id, parseErr := strconv.ParseUint(idStr, 10, 64)
	if parseErr != nil || id == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}

	var req struct {
		ParentID uint64 `json:"parent_id"`
		Content  string `json:"content"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return httpx.Error(ctx, 422, "回复内容不能为空")
	}
	if len(content) > 5000 {
		return httpx.Error(ctx, 422, "回复内容过长（最多 5000 字）")
	}

	throttleKey := "reply:throttle:" + strconv.FormatUint(identity.ID, 10)
	if parseCount(facades.Cache().Get(throttleKey)) > 0 {
		return httpx.Error(ctx, 429, "回复太快了，请稍等几秒再试")
	}
	_ = facades.Cache().Put(throttleKey, 1, replyThrottleWindow)

	var statuses []string
	if err := facades.Orm().Query().Table("topics").
		Where("id = ? AND deleted_at IS NULL", id).
		Pluck("status", &statuses); err != nil || len(statuses) == 0 {
		return httpx.Error(ctx, 404, "topic not found")
	}
	if statuses[0] != "open" {
		return httpx.Error(ctx, 422, "帖子已关闭，无法回复")
	}

	reply, err := forumservices.ReplyToTopic(id, identity.ID, req.ParentID, content)
	if err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}

	// Daily-capped earn reward (best-effort; never blocks the reply).
	_, _ = settingsservices.AwardDaily(identity.ID, settingsservices.PointReasonReply,
		settingsservices.EarnReply(), settingsservices.ReplyDailyCap())

	// Fan out @mentions to the mentioned users' notification feeds.
	notificationservices.NotifyMentions(identity.ID, identity.Name, content, "topic", id, content)

	return ctx.Response().Success().Json(http.Json{
		"id":         reply.ID,
		"user_id":    reply.UserID,
		"parent_id":  reply.ParentID,
		"content":    reply.Content,
		"floor":      reply.Floor,
		"like_count": reply.LikeCount,
		"created_at": reply.CreatedAt,
	})
}

func parseCount(raw any) int64 {
	s := fmt.Sprint(raw)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func clampPageV1(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 1
	}
	return n
}

func clampPerPageV1(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 20
	}
	if n > 50 {
		return 50
	}
	return n
}

func envelope(ctx http.Context, items []map[string]any, total, page, perPage int) http.Response {
	totalPages := (total + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}
	return ctx.Response().Success().Json(http.Json{
		"items":      items,
		"total":      total,
		"page":       page,
		"perPage":    perPage,
		"totalPages": totalPages,
	})
}
