package http

import (
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	authhttp "reflexcms/backend/modules/auth/http/controllers"
	notificationservices "reflexcms/backend/modules/notification/services"
	settingsservices "reflexcms/backend/modules/settings/services"
)

// gateOr404 returns a 404 response if the section is disabled by site.mode.
func gateOr404(ctx http.Context, section string) *http.Response {
	if !settingsservices.SectionEnabled(section) {
		resp := httpx.Error(ctx, 404, "Not Found")
		return &resp
	}
	return nil
}

// PublicCommentController serves the /api/v1 comment endpoints: reads are
// anonymous, writes require a session (identity always comes from the token,
// never from the payload).
type PublicCommentController struct{}

func NewPublicCommentController() *PublicCommentController {
	return &PublicCommentController{}
}

// List: GET /api/v1/articles/:slug/comments — approved comments, paginated.
func (r *PublicCommentController) List(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "comments"); resp != nil {
		return *resp
	}

	slug := ctx.Request().Route("slug")

	var articleID int64
	if err := facades.Orm().Query().Table("articles").
		Where("slug = ? AND deleted_at IS NULL", slug).
		Pluck("id", &articleID); err != nil || articleID == 0 {
		return ctx.Response().Success().Json(http.Json{"items": []any{}, "total": 0})
	}

	page, perPage := 1, 20
	if n, err := strconv.Atoi(ctx.Request().Query("page", "1")); err == nil && n > 0 {
		page = n
	}
	if n, err := strconv.Atoi(ctx.Request().Query("perPage", "20")); err == nil && n > 0 && n <= 50 {
		perPage = n
	}

	query := facades.Orm().Query().Table("comments").
		Where("article_id = ? AND status = 'approved' AND deleted_at IS NULL", articleID)

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	rows := []map[string]any{}
	if err := query.
		Select("id", "user_id", "parent_id", "content", "created_at").
		OrderBy("created_at", "asc").
		Offset((page - 1) * perPage).Limit(perPage).
		Get(&rows); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	totalPages := (int(total) + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}
	return ctx.Response().Success().Json(http.Json{
		"items": rows, "total": total, "page": page,
		"perPage": perPage, "totalPages": totalPages,
	})
}

// Submit: POST /api/v1/articles/:slug/comments — creates a pending comment
// attributed to the session user.
func (r *PublicCommentController) Submit(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "comments"); resp != nil {
		return *resp
	}

	identity, ok := authhttp.RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, 401, "请先登录后再评论")
	}

	slug := ctx.Request().Route("slug")
	var articleID int64
	if err := facades.Orm().Query().Table("articles").
		Where("slug = ? AND deleted_at IS NULL", slug).
		Pluck("id", &articleID); err != nil || articleID == 0 {
		return httpx.Error(ctx, 404, "article not found")
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "content is required")
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return httpx.Error(ctx, 422, "评论内容不能为空")
	}
	if len(content) > 2000 {
		return httpx.Error(ctx, 422, "评论内容过长（最多 2000 字）")
	}

	if _, err := facades.Orm().Query().Exec(`
		INSERT INTO comments (article_id, user_id, parent_id, content, status, created_at, updated_at)
		VALUES (?, ?, 0, ?, 'pending', NOW(), NOW())
		RETURNING id
	`, articleID, identity.ID, content); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	// Fan out @mentions (link points at the article; excerpt carries context).
	notificationservices.NotifyMentions(identity.ID, identity.Name, content, "article", uint64(articleID), content)

	return ctx.Response().Success().Json(http.Json{
		"message": "评论已提交，等待审核",
	})
}

// Ensure settingsservices is used (gateOr404 references it).
var _ = settingsservices.SectionEnabled
