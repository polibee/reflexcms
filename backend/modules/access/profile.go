package access

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	authhttp "reflexcms/backend/modules/auth/http/controllers"
	accessmodels "reflexcms/backend/modules/access/models"
	settingsservices "reflexcms/backend/modules/settings/services"
)

/* Public /api/v1 profile endpoints: user homepages (概览/主题帖/回复/收藏)
 * and self-service profile & security. All reads are anonymous; all writes
 * require a session and derive identity from the bearer token. */

func clampPage(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 1
	}
	return n
}

func clampPerPage(raw string, max int) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 20
	}
	if n > max {
		return max
	}
	return n
}

func envelope(ctx http.Context, items []map[string]any, total, page, perPage int) http.Response {
	totalPages := (total + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}
	return ctx.Response().Success().Json(http.Json{
		"items": items, "total": total, "page": page,
		"perPage": perPage, "totalPages": totalPages,
	})
}

func parseIDParam(ctx http.Context) (uint64, bool) {
	id, err := strconv.ParseUint(ctx.Request().Route("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

// ShowUser: GET /api/v1/users/{id} — public profile with content counts.
func ShowUser(ctx http.Context) http.Response {
	id, ok := parseIDParam(ctx)
	if !ok {
		return httpx.Error(ctx, 404, "user not found")
	}

	var rows []map[string]any
	if err := facades.Orm().Query().Table("users").
		Where("id = ? AND status = 'active' AND deleted_at IS NULL", id).
		Select("id", "username", "avatar", "bio", "website", "profile_card", "created_at").
		Get(&rows); err != nil || len(rows) == 0 {
		return httpx.Error(ctx, 404, "user not found")
	}

	profile := rows[0]

	topicCount, _ := facades.Orm().Query().Table("topics").
		Where("user_id = ? AND deleted_at IS NULL", id).Count()
	replyCount, _ := facades.Orm().Query().Table("replies").
		Where("user_id = ? AND deleted_at IS NULL", id).Count()
	favCount, _ := facades.Orm().Query().Table("favorites").
		Where("user_id = ? AND fav_type = ?", id, "topic").Count()

	profile["counts"] = http.Json{
		"topics": topicCount, "replies": replyCount, "favorites": favCount,
	}
	return ctx.Response().Success().Json(profile)
}

// UserTopics: GET /api/v1/users/{id}/topics — paginated topics by the user.
func UserTopics(ctx http.Context) http.Response {
	id, ok := parseIDParam(ctx)
	if !ok {
		return httpx.Error(ctx, 404, "user not found")
	}
	page := clampPage(ctx.Request().Query("page", "1"))
	perPage := clampPerPage(ctx.Request().Query("perPage", "20"), 50)

	query := facades.Orm().Query().Table("topics").
		Where("user_id = ? AND status = 'open' AND deleted_at IS NULL", id)

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.Select("id", "title", "forum_category_id", "view_count", "reply_count", "created_at").
		OrderBy("created_at", "desc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return envelope(ctx, items, int(total), page, perPage)
}

// UserReplies: GET /api/v1/users/{id}/replies — paginated replies with the
// parent topic title so the profile page can link back.
func UserReplies(ctx http.Context) http.Response {
	id, ok := parseIDParam(ctx)
	if !ok {
		return httpx.Error(ctx, 404, "user not found")
	}
	page := clampPage(ctx.Request().Query("page", "1"))
	perPage := clampPerPage(ctx.Request().Query("perPage", "20"), 50)

	query := facades.Orm().Query().Table("replies as r").
		Join("JOIN topics t ON t.id = r.topic_id AND t.deleted_at IS NULL").
		Where("r.user_id = ? AND r.deleted_at IS NULL", id)

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.Select("r.id", "r.topic_id", "t.title AS topic_title", "r.content", "r.floor", "r.created_at").
		OrderBy("r.created_at", "desc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return envelope(ctx, items, int(total), page, perPage)
}

// UserFavorites: GET /api/v1/users/{id}/favorites — paginated favorited
// topics (newest bookmark first).
func UserFavorites(ctx http.Context) http.Response {
	id, ok := parseIDParam(ctx)
	if !ok {
		return httpx.Error(ctx, 404, "user not found")
	}
	page := clampPage(ctx.Request().Query("page", "1"))
	perPage := clampPerPage(ctx.Request().Query("perPage", "20"), 50)

	aliveFilter := `(
		(f.fav_type = 'topic' AND EXISTS (SELECT 1 FROM topics t WHERE t.id = f.fav_id AND t.deleted_at IS NULL))
		OR (f.fav_type = 'article' AND EXISTS (SELECT 1 FROM articles a WHERE a.id = f.fav_id AND a.deleted_at IS NULL))
	)`

	query := facades.Orm().Query().Table("favorites as f").
		Join("LEFT JOIN topics t ON t.id = f.fav_id AND f.fav_type = 'topic' AND t.deleted_at IS NULL").
		Join("LEFT JOIN articles a ON a.id = f.fav_id AND f.fav_type = 'article' AND a.deleted_at IS NULL").
		Where("f.user_id = ?", id).
		Where(aliveFilter)

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.Select(`f.id, f.fav_type, f.fav_id, f.created_at AS favorited_at,
		COALESCE(t.title, a.title) AS title,
		COALESCE(t.reply_count, 0) AS reply_count,
		a.slug AS article_slug`).
		OrderBy("f.created_at", "desc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return envelope(ctx, items, int(total), page, perPage)
}

// MyNotifications: GET /api/v1/me/notifications — the session user's feed
// (mentions, moderation results, reply notices), newest first.
func MyNotifications(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}
	page := clampPage(ctx.Request().Query("page", "1"))
	perPage := clampPerPage(ctx.Request().Query("perPage", "20"), 50)

	query := facades.Orm().Query().Table("notifications").
		Where("user_id = ?", identity.ID)

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.
		Select("id", "type", "data", "read_at", "created_at").
		OrderBy("created_at", "desc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return envelope(ctx, items, int(total), page, perPage)
}

// ReadAllNotifications: POST /api/v1/me/notifications/read-all.
func ReadAllNotifications(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}
	if _, err := facades.Orm().Query().Exec(`
		UPDATE notifications SET read_at = NOW()
		WHERE user_id = ? AND read_at IS NULL
	`, identity.ID); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"message": "全部已读"})
}

/* ---------------- self-service ---------------- */

// MyProfile: GET /api/v1/me/profile — the session user's editable profile.
func MyProfile(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var rows []map[string]any
	if err := facades.Orm().Query().Table("users").
		Where("id = ?", identity.ID).
		Select("id", "username", "email", "avatar", "bio", "website", "signature", "profile_card", "created_at").
		Get(&rows); err != nil || len(rows) == 0 {
		return httpx.Error(ctx, 404, "user not found")
	}
	return ctx.Response().Success().Json(rows[0])
}

// MyStats: GET /api/v1/me/stats — counters for the sidebar user-center card
// (topics, forum replies, article comments, favorites, mention
// notifications). Mentions have no producer yet; the key stays so the card
// shape is stable when it lands.
func MyStats(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	uid := identity.ID
	topicCount, _ := facades.Orm().Query().Table("topics").
		Where("user_id = ? AND deleted_at IS NULL", uid).Count()
	replyCount, _ := facades.Orm().Query().Table("replies").
		Where("user_id = ? AND deleted_at IS NULL", uid).Count()
	commentCount, _ := facades.Orm().Query().Table("comments").
		Where("user_id = ? AND deleted_at IS NULL", uid).Count()
	favCount, _ := facades.Orm().Query().Table("favorites").
		Where("user_id = ? AND fav_type = ?", uid, "topic").Count()
	mentionCount, _ := facades.Orm().Query().Table("notifications").
		Where("user_id = ? AND type = ?", uid, "mention").Count()
	unreadDM, _ := facades.Orm().Query().Table("messages").
		Where("recipient_id = ? AND read_at IS NULL", uid).Count()

	return ctx.Response().Success().Json(http.Json{
		"topics":          topicCount,
		"replies":         replyCount,
		"comments":        commentCount,
		"favorites":       favCount,
		"mentions":        mentionCount,
		"messages_unread": unreadDM,
	})
}

/* ---------------- points / level / daily tasks ---------------- */

// MyPoints: GET /api/v1/me/points — balance, level progress and today's
// task quotas for the user-center card.
func MyPoints(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	balance := settingsservices.GetBalance(identity.ID)
	total := settingsservices.TotalEarned(identity.ID)
	level, curTh, nextTh := settingsservices.LevelInfo(total)

	percent := 100
	if nextTh > curTh {
		percent = (total - curTh) * 100 / (nextTh - curTh)
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}
	}

	topicEarned := settingsservices.DailyEarned(identity.ID, settingsservices.PointReasonTopic)
	replyEarned := settingsservices.DailyEarned(identity.ID, settingsservices.PointReasonReply)
	signinDone := settingsservices.HasReasonToday(identity.ID, settingsservices.PointReasonSignin)

	return ctx.Response().Success().Json(http.Json{
		"currency":     settingsservices.CurrencyName(),
		"balance":      balance,
		"total_earned": total,
		"level": http.Json{
			"current":    level,
			"next":       level + 1,
			"percent":    percent,
			"cur_floor":  curTh,
			"next_floor": nextTh,
		},
		"daily": http.Json{
			"topic":  http.Json{"earned": topicEarned, "cap": settingsservices.TopicDailyCap(), "reward": settingsservices.EarnTopic()},
			"reply":  http.Json{"earned": replyEarned, "cap": settingsservices.ReplyDailyCap(), "reward": settingsservices.EarnReply()},
			"signin": http.Json{"done": signinDone, "reward": settingsservices.EarnSignin()},
		},
	})
}

// SignIn: POST /api/v1/me/signin — once per day; grants the configured
// reward and reports today's task state back.
func SignIn(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	reason := settingsservices.PointReasonSignin
	if settingsservices.HasReasonToday(identity.ID, reason) {
		return httpx.Error(ctx, 422, "今天已经签到过了，明天再来吧")
	}

	reward := settingsservices.EarnSignin()
	if err := settingsservices.GrantPoints(identity.ID, reward, reason); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	return ctx.Response().Success().Json(http.Json{
		"message": fmt.Sprintf("签到成功，获得 %d %s", reward, settingsservices.CurrencyName()),
		"reward":  reward,
	})
}

/* ---------------- direct messages ---------------- */

// SendMessage: POST /api/v1/me/messages — {to_username, body}.
func SendMessage(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var req struct {
		ToUsername string `json:"to_username"`
		Body       string `json:"body"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}
	to := strings.TrimSpace(req.ToUsername)
	body := strings.TrimSpace(req.Body)
	if to == "" || body == "" {
		return httpx.Error(ctx, 422, "收件人和内容不能为空")
	}
	if len([]rune(body)) > 2000 {
		return httpx.Error(ctx, 422, "私信内容过长（最多 2000 字）")
	}

	var recipientIDs []int64
	if err := facades.Orm().Query().Table("users").
		Where("username = ? AND status = 'active' AND deleted_at IS NULL", to).
		Pluck("id", &recipientIDs); err != nil || len(recipientIDs) == 0 {
		return httpx.Error(ctx, 404, "收件人不存在")
	}
	recipientID := uint64(recipientIDs[0])
	if recipientID == identity.ID {
		return httpx.Error(ctx, 422, "不能给自己发私信")
	}

	if _, err := facades.Orm().Query().Exec(`
		INSERT INTO messages (sender_id, recipient_id, body, created_at)
		VALUES (?, ?, ?, NOW())
	`, identity.ID, recipientID, body); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	return ctx.Response().Success().Json(http.Json{"message": "私信已发送"})
}

// MyMessages: GET /api/v1/me/messages?box=inbox|sent — paginated.
func MyMessages(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	box := ctx.Request().Query("box", "inbox")
	page := clampPage(ctx.Request().Query("page", "1"))
	perPage := clampPerPage(ctx.Request().Query("perPage", "20"), 50)

	// counterpart label: inbox shows sender, sent shows recipient.
	otherJoin := "sender"
	otherCol := "m.sender_id"
	otherKey := "from_name"
	where := "m.recipient_id = ?"
	if box == "sent" {
		otherJoin = "recipient"
		otherCol = "m.recipient_id"
		otherKey = "to_name"
		where = "m.sender_id = ?"
	}

	query := facades.Orm().Query().Table("messages as m").
		Join("JOIN users u ON u.id = "+otherCol).
		Where(where, identity.ID)

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.
		Select("m.id", "m.sender_id", "m.recipient_id", "m.body", "m.read_at", "m.created_at", "u.username AS "+otherKey).
		OrderBy("m.created_at", "desc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	_ = otherJoin
	return envelope(ctx, items, int(total), page, perPage)
}

// UnreadMessages: GET /api/v1/me/messages/unread-count — badge number for
// the user-center card; anonymous callers get zero.
func UnreadMessages(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return ctx.Response().Success().Json(http.Json{"unread": 0})
	}
	var count int64
	count, _ = facades.Orm().Query().Table("messages").
		Where("recipient_id = ? AND read_at IS NULL", identity.ID).Count()
	return ctx.Response().Success().Json(http.Json{"unread": count})
}

// MarkMessageRead: POST /api/v1/me/messages/{id}/read — only the recipient
// may mark their own message.
func MarkMessageRead(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}
	id, valid := parseIDParam(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "message not found")
	}
	if _, err := facades.Orm().Query().Exec(`
		UPDATE messages SET read_at = NOW()
		WHERE id = ? AND recipient_id = ? AND read_at IS NULL
	`, id, identity.ID); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"message": "已读"})
}

/* ---------------- blocking ---------------- */

// ToggleBlock: POST /api/v1/me/blocked/toggle — {user_id}. Idempotent flip;
// blocked users' replies disappear from the caller's feeds.
func ToggleBlock(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var req struct {
		UserID uint64 `json:"user_id"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.UserID == 0 {
		return httpx.Error(ctx, 422, "user_id is required")
	}
	if req.UserID == identity.ID {
		return httpx.Error(ctx, 422, "不能拉黑自己")
	}

	var target int64
	target, _ = facades.Orm().Query().Table("users").
		Where("id = ? AND deleted_at IS NULL", req.UserID).Count()
	if target == 0 {
		return httpx.Error(ctx, 404, "user not found")
	}

	var current int64
	current, _ = facades.Orm().Query().Table("blocked_users").
		Where("user_id = ? AND blocked_id = ?", identity.ID, req.UserID).Count()

	if current > 0 {
		if _, err := facades.Orm().Query().Exec(`
			DELETE FROM blocked_users WHERE user_id = ? AND blocked_id = ?
		`, identity.ID, req.UserID); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"blocked": false, "message": "已取消拉黑"})
	}

	if _, err := facades.Orm().Query().Exec(`
		INSERT INTO blocked_users (user_id, blocked_id, created_at)
		VALUES (?, ?, NOW())
		ON CONFLICT (user_id, blocked_id) DO NOTHING
	`, identity.ID, req.UserID); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"blocked": true, "message": "已拉黑，将不再看到该用户的回复"})
}

// MyBlocked: GET /api/v1/me/blocked — the block list with usernames.
func MyBlocked(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	items := []map[string]any{}
	if err := facades.Orm().Query().Table("blocked_users as b").
		Join("JOIN users u ON u.id = b.blocked_id").
		Where("b.user_id = ?", identity.ID).
		Select("b.id", "b.blocked_id", "u.username", "b.created_at").
		OrderBy("b.created_at", "desc").
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	if items == nil {
		items = []map[string]any{}
	}
	return ctx.Response().Success().Json(http.Json{"items": items})
}

// BlockedIDs returns the caller's block list; empty when anonymous.
func BlockedIDs(ctx http.Context) []uint64 {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return nil
	}
	var ids []int64
	_ = facades.Orm().Query().Table("blocked_users").
		Where("user_id = ?", identity.ID).
		Pluck("blocked_id", &ids)
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		out = append(out, uint64(id))
	}
	return out
}

// UpdateMyProfileCard: PUT /api/v1/me/profile-card — save the Markdown
// self-description card shown on the public user homepage. Empty content
// hides the card entirely.
func UpdateMyProfileCard(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}
	content := strings.TrimSpace(req.Content)
	if len([]rune(content)) > 10000 {
		return httpx.Error(ctx, 422, "卡片内容过长（最多 10000 字）")
	}

	if _, err := facades.Orm().Query().Table("users").
		Where("id = ?", identity.ID).
		Update(map[string]any{"profile_card": content}); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	msg := "卡片已保存"
	if content == "" {
		msg = "卡片已清空"
	}
	return ctx.Response().Success().Json(http.Json{"message": msg})
}

// dataNormalize trims a string field; present=false means the caller did
// not send it (absent key), letting updates distinguish "clear" vs "skip".
func dataNormalize(s string) (string, bool) {
	return strings.TrimSpace(s), strings.TrimSpace(s) != "" || s != ""
}

// UpdateMyProfile: PUT /api/v1/me/profile — edit username/bio/website/avatar.
func UpdateMyProfile(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var req struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Website   string `json:"website"`
		Avatar    string `json:"avatar"`
		Signature string `json:"signature"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}

	updates := map[string]any{}
	if name := strings.TrimSpace(req.Username); name != "" {
		if len([]rune(name)) > 32 {
			return httpx.Error(ctx, 422, "用户名过长（最多 32 字）")
		}
		var count int64
		count, _ = facades.Orm().Query().Table("users").
			Where("username = ? AND id <> ?", name, identity.ID).Count()
		if count > 0 {
			return httpx.Error(ctx, 422, "用户名已被占用")
		}
		updates["username"] = name
	}
	// Signature: short markdown, rendered client-side through the escaping
	// renderer (no HTML/images/scripts survive); plain text + links only.
	if sig, present := dataNormalize(req.Signature); present {
		if len([]rune(sig)) > 200 {
			return httpx.Error(ctx, 422, "签名过长（最多 200 字）")
		}
		updates["signature"] = sig
	}
	if req.Bio != "" || req.Website != "" || req.Avatar != "" || req.Signature != "" {
		if len([]rune(req.Bio)) > 500 {
			return httpx.Error(ctx, 422, "简介过长（最多 500 字）")
		}
		updates["bio"] = req.Bio
		updates["website"] = strings.TrimSpace(req.Website)
		updates["avatar"] = strings.TrimSpace(req.Avatar)
	}
	if len(updates) == 0 {
		return httpx.Error(ctx, 422, "没有需要更新的字段")
	}

	if _, err := facades.Orm().Query().Table("users").
		Where("id = ?", identity.ID).Update(updates); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"message": "资料已更新"})
}

// ChangeMyPassword: PUT /api/v1/me/password — verify current, set new.
func ChangeMyPassword(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var req struct {
		Current string `json:"current_password"`
		New     string `json:"new_password"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}
	if len(req.New) < 8 {
		return httpx.Error(ctx, 422, "新密码至少 8 位")
	}

	var user accessmodels.User
	if err := facades.Orm().Query().Where("id = ?", identity.ID).FirstOrFail(&user); err != nil {
		return httpx.Error(ctx, 404, "user not found")
	}
	if !facades.Hash().Check(req.Current, user.Password) {
		return httpx.Error(ctx, 422, "当前密码不正确")
	}

	hashed, err := facades.Hash().Make(req.New)
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	if _, err := facades.Orm().Query().Table("users").
		Where("id = ?", identity.ID).Update(map[string]any{"password": hashed}); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"message": "密码已修改"})
}

/* ---------------- favorites ---------------- */

// ToggleFavorite: POST /api/v1/favorites/toggle — idempotent bookmark flip
// for the session user. Body: {type: "topic", id: 123}.
func ToggleFavorite(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录")
	}

	var req struct {
		Type string `json:"type"`
		ID   uint64 `json:"id"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.ID == 0 {
		return httpx.Error(ctx, 422, "id is required")
	}
	if req.Type == "" {
		req.Type = "topic"
	}
	if req.Type != "topic" && req.Type != "article" {
		return httpx.Error(ctx, 422, "unsupported favorite type")
	}

	switch req.Type {
	case "topic":
		var exists int64
		exists, _ = facades.Orm().Query().Table("topics").
			Where("id = ? AND deleted_at IS NULL", req.ID).Count()
		if exists == 0 {
			return httpx.Error(ctx, 404, "topic not found")
		}
	case "article":
		var exists int64
		exists, _ = facades.Orm().Query().Table("articles").
			Where("id = ? AND deleted_at IS NULL", req.ID).Count()
		if exists == 0 {
			return httpx.Error(ctx, 404, "article not found")
		}
	}

	var current int64
	current, _ = facades.Orm().Query().Table("favorites").
		Where("user_id = ? AND fav_type = ? AND fav_id = ?", identity.ID, req.Type, req.ID).
		Count()

	if current > 0 {
		_, err := facades.Orm().Query().Table("favorites").
			Where("user_id = ? AND fav_type = ? AND fav_id = ?", identity.ID, req.Type, req.ID).
			Delete()
		if err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"favorited": false})
	}

	if _, err := facades.Orm().Query().Exec(`
		INSERT INTO favorites (user_id, fav_type, fav_id, created_at)
		VALUES (?, ?, ?, NOW())
		ON CONFLICT (user_id, fav_type, fav_id) DO NOTHING
	`, identity.ID, req.Type, req.ID); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{"favorited": true})
}

// MyFavoriteState: GET /api/v1/favorites/state?type=topic&id=123 — whether
// the session user has favorited the target (drives the button state).
func MyFavoriteState(ctx http.Context) http.Response {
	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return ctx.Response().Success().Json(http.Json{"favorited": false})
	}
	favType := ctx.Request().Query("type", "topic")
	id, err := strconv.ParseUint(ctx.Request().Query("id", "0"), 10, 64)
	if err != nil || id == 0 {
		return ctx.Response().Success().Json(http.Json{"favorited": false})
	}
	var count int64
	count, _ = facades.Orm().Query().Table("favorites").
		Where("user_id = ? AND fav_type = ? AND fav_id = ?", identity.ID, favType, id).Count()
	return ctx.Response().Success().Json(http.Json{"favorited": count > 0})
}
