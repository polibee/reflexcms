// Package pages is a system-level module: admin-authored static pages
// (隐私政策、服务条款、关于我们…) served at /api/v1/pages/{slug}. It sits at
// the platform layer because pages are consumed by every product
// composition (CMS and forum alike) rather than belonging to either.
package pages

import (
	"strings"

	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	"reflexcms/backend/app/support/slug"
	"reflexcms/backend/database/migrations"
	adminhub "reflexcms/backend/modules/adminhub"
	pagesmodels "reflexcms/backend/modules/pages/models"
)

// Module owns the static-pages feature.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "pages" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.PagesMigrations()
}

func (m *Module) Routes() {
	adminhub.Register(&adminhub.Spec{
		Name:       "pages",
		Model:      pagesmodels.Page{},
		Searchable: []string{"title", "slug"},
		Sortable:   []string{"id", "title", "slug", "status", "created_at"},
		Fillable:   []string{"title", "slug", "content", "status"},
		BeforeSave: func(data map[string]any) error {
			slugVal, _ := data["slug"].(string)
			if strings.TrimSpace(slugVal) != "" {
				return nil
			}
			title, _ := data["title"].(string)
			if title == "" {
				return nil
			}
			unique, err := slug.EnsureUnique(slug.Slugify(title), 191, func(candidate string) (bool, error) {
				var count int64
				count, err := facades.Orm().Query().
					Model(pagesmodels.Page{}).
					Where("slug = ?", candidate).
					Count()
				return count > 0, err
			})
			if err != nil {
				return err
			}
			data["slug"] = unique
			return nil
		},
		PermissionPrefix: "pages",
	})

	registerPublicRoutes()
}

func (m *Module) Boot() {}

// registerPublicRoutes exposes published pages anonymously.
func registerPublicRoutes() {
	r := facades.Route()
	r.Get("/api/v1/pages/{slug}", func(ctx http.Context) http.Response {
		slug := ctx.Request().Route("slug")
		var rows []map[string]any
		if err := facades.Orm().Query().Table("pages").
			Where("slug = ? AND status = 'published'", slug).
			Select("id", "title", "slug", "content", "updated_at").
			Get(&rows); err != nil || len(rows) == 0 {
			return httpx.Error(ctx, 404, "page not found")
		}
		return ctx.Response().Success().Json(rows[0])
	})
}
