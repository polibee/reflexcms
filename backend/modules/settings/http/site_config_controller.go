package http

import (
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	settingsservices "reflexcms/backend/modules/settings/services"
)

type SiteConfigController struct{}

func NewSiteConfigController() *SiteConfigController { return &SiteConfigController{} }

// Show: GET /api/v1/site/config — the first call any public frontend makes.
// Anonymous by design; exposes only non-sensitive presentation config plus
// which content sections the current site.mode exposes (plan §5.4).
func (r *SiteConfigController) Show(ctx http.Context) http.Response {
	mode := settingsservices.Mode()

	home := settingsservices.Get("site.home")
	if home != "cms" && home != "forum" {
		home = "both"
	}

	return ctx.Response().Success().Json(http.Json{
		"name": facades.Config().GetString("app.name", "ReflexCMS"),
		"mode": mode,
		"home": home,
		"seo": map[string]any{
			"description": settingsservices.Get("site.seo.description"),
			"icp":         settingsservices.Get("site.icp"),
		},
		"sections": sectionsFor(mode),
	})
}

func sectionsFor(mode string) []string {
	switch mode {
	case settingsservices.ModeBlog:
		return []string{"articles", "comments"}
	case settingsservices.ModeForum:
		return []string{"forums", "topics", "replies"}
	default:
		return []string{"articles", "comments", "forums", "topics", "replies"}
	}
}
