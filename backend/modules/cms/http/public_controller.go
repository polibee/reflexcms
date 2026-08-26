package http

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	settingsservices "reflexcms/backend/modules/settings/services"
)

// PublicController serves anonymous /api/v1 read endpoints for the CMS
// domain. Every handler checks site.mode gating before touching data.
type PublicController struct{}

func NewPublicController() *PublicController { return &PublicController{} }

func gateOr404(ctx http.Context, section string) *http.Response {
	if !settingsservices.SectionEnabled(section) {
		resp := httpx.Error(ctx, 404, "Not Found")
		return &resp
	}
	return nil
}

// ListArticles: GET /api/v1/articles — published only, paginated.
func (r *PublicController) ListArticles(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "articles"); resp != nil {
		return *resp
	}

	page := clampPage(ctx.Request().Query("page", "1"))
	perPage := clampPerPage(ctx.Request().Query("perPage", "20"))

	query := facades.Orm().Query().Table("articles").
		Where("status = ? AND deleted_at IS NULL", "published").
		Select("id", "title", "slug", "summary", "cover", "is_pinned", "published_at", "view_count", "comment_count")

	if term := ctx.Request().Query("q", ""); term != "" {
		query = query.Where("search_vector @@ websearch_to_tsquery('simple', ?)", term)
	}
	if catID := ctx.Request().Query("category_id", ""); catID != "" {
		query = query.Where("category_id = ?", catID)
	}

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := []map[string]any{}
	if err := query.
		OrderBy("is_pinned", "desc").
		OrderBy("published_at", "desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Get(&items); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	return envelope(ctx, items, int(total), page, perPage)
}

// ShowArticle: GET /api/v1/articles/:slug — published article detail.
// Also increments view_count via Redis buffer.
func (r *PublicController) ShowArticle(ctx http.Context) http.Response {
	if resp := gateOr404(ctx, "articles"); resp != nil {
		return *resp
	}

	slug := ctx.Request().Route("slug")
	var rows []map[string]any
	if err := facades.Orm().Query().Table("articles").
		Where("slug = ? AND status = ? AND deleted_at IS NULL", slug, "published").
		Select("id", "title", "slug", "content", "summary", "cover",
			"is_pinned", "published_at", "view_count", "comment_count", "created_at").
		Limit(1).
		Get(&rows); err != nil || len(rows) == 0 {
		return httpx.Error(ctx, 404, "article not found")
	}

	// Buffered view count increment (flushed by scheduler every 5 min).
	facades.Cache().Increment("views:article:"+slug)

	return ctx.Response().Success().Json(rows[0])
}

/* ---------------- pagination helpers ---------------- */

func clampPage(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 1
	}
	return n
}

func clampPerPage(raw string) int {
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
