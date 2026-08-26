package http

import (
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	adminhub "reflexcms/backend/modules/adminhub"
	auditservices "reflexcms/backend/modules/audit/services"
	cmsmodels "reflexcms/backend/modules/cms/models"
	cmsservices "reflexcms/backend/modules/cms/services"
)

var acceptedTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// ArticleActionController serves the dedicated lifecycle routes:
//
//	POST /api/admin/articles/{id}/publish|schedule|unpublish|archive|restore
//
// publish/schedule accept an optional {"at": "<time>"} body — a future time
// turns the article into a scheduled release picked up by the scheduler.
// Status changes never flow through the generic update payload so every
// transition passes through the state machine.
type ArticleActionController struct{}

func NewArticleActionController() *ArticleActionController {
	return &ArticleActionController{}
}

func (r *ArticleActionController) Publish(ctx http.Context) http.Response {
	return r.run(ctx, "publish")
}

func (r *ArticleActionController) Schedule(ctx http.Context) http.Response {
	return r.run(ctx, "schedule")
}

func (r *ArticleActionController) Unpublish(ctx http.Context) http.Response {
	return r.run(ctx, "unpublish")
}

func (r *ArticleActionController) Archive(ctx http.Context) http.Response {
	return r.run(ctx, "archive")
}

func (r *ArticleActionController) Restore(ctx http.Context) http.Response {
	return r.run(ctx, "restore")
}

// RestoreFromTrash revives a soft-deleted article. This static route shadows
// the gateway's generic restore, so trashed-row revival lives here too.
func (r *ArticleActionController) RestoreFromTrash(ctx http.Context) http.Response {
	if resp := adminhub.CheckPermission(ctx, "articles.delete"); resp != nil {
		return *resp
	}
	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "article not found")
	}

	if _, err := facades.Orm().Query().Exec(
		"UPDATE articles SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL", id,
	); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	var article cmsmodels.Article
	if err := facades.Orm().Query().FindOrFail(&article, id); err != nil {
		return httpx.Error(ctx, 404, "article not found")
	}
	return ctx.Response().Success().Json(http.Json{"id": article.ID, "status": article.Status})
}

// Pin / Unpin: WordPress-style sticky articles.
func (r *ArticleActionController) Pin(ctx http.Context) http.Response {
	return r.setPinned(ctx, true)
}

func (r *ArticleActionController) Unpin(ctx http.Context) http.Response {
	return r.setPinned(ctx, false)
}

func (r *ArticleActionController) setPinned(ctx http.Context, value bool) http.Response {
	if resp := adminhub.CheckPermission(ctx, "articles.edit"); resp != nil {
		return *resp
	}
	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "article not found")
	}

	var article cmsmodels.Article
	if err := facades.Orm().Query().FindOrFail(&article, id); err != nil {
		return httpx.Error(ctx, 404, "article not found")
	}
	article.IsPinned = value
	if err := facades.Orm().Query().Save(&article); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	auditservices.Record(ctx, adminhub.CurrentUserID(ctx), pinAction(value), "articles", article.ID, http.StatusOK, nil)

	return ctx.Response().Success().Json(http.Json{"id": article.ID, "is_pinned": article.IsPinned})
}

func pinAction(value bool) string {
	if value {
		return "pin"
	}
	return "unpin"
}

func (r *ArticleActionController) run(ctx http.Context, action string) http.Response {
	required := "articles.edit" // unpublish/archive/restore are editorial edits
	if action == "publish" || action == "schedule" {
		required = "articles.publish"
	}
	if resp := adminhub.CheckPermission(ctx, required); resp != nil {
		return *resp
	}

	id, valid := adminhub.ParseRowID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, "article not found")
	}

	var article cmsmodels.Article
	if err := facades.Orm().Query().FindOrFail(&article, id); err != nil {
		return httpx.Error(ctx, 404, "article not found")
	}

	to, ok := cmsservices.TransitionAction(action)
	if !ok {
		return httpx.Error(ctx, 404, "unknown action \""+action+"\"")
	}

	// Optional scheduled release time.
	body := map[string]any{}
	_ = ctx.Request().Bind(&body)
	at := parseFlexibleTime(body["at"])
	if (action == "publish" || action == "schedule") && at != nil &&
		at.After(time.Now()) {
		article.PublishedAt = toCarbon(at)
		to = cmsmodels.ArticleScheduled
	}
	if action == "schedule" && at == nil {
		return httpx.Error(ctx, 422, "\"at\" is required for scheduling")
	}

	if err := cmsservices.Transition(&article, to); err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}

	auditservices.Record(ctx, adminhub.CurrentUserID(ctx), action, "articles", article.ID, http.StatusOK, map[string]any{
		"status":       article.Status,
		"published_at": article.PublishedAt,
	})

	return ctx.Response().Success().Json(http.Json{
		"id":           article.ID,
		"status":       article.Status,
		"published_at": article.PublishedAt,
	})
}

func parseFlexibleTime(raw any) *time.Time {
	s, ok := raw.(string)
	if !ok {
		return nil
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	for _, layout := range acceptedTimeLayouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return &t
		}
	}
	return nil
}

func toCarbon(t *time.Time) *carbon.DateTime {
	return carbon.NewDateTime(carbon.FromStdTime(*t))
}
