package cms

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"

	"reflexcms/backend/app/facades"
	adminhub "reflexcms/backend/modules/adminhub"
	cmshttp "reflexcms/backend/modules/cms/http"
	cmsmodels "reflexcms/backend/modules/cms/models"
	cmsservices "reflexcms/backend/modules/cms/services"
)

// Register exposes the CMS domain: three gateway resources plus dedicated
// lifecycle action routes for articles.
func Register() {
	adminhub.Register(categorySpec())
	adminhub.Register(tagSpec())
	adminhub.Register(articleSpec())
	registerActionRoutes()
}

/* ---------------- categories ---------------- */

func categorySpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "categories",
		Model:      cmsmodels.ArticleCategory{},
		Searchable: []string{"name", "slug"},
		Sortable:   []string{"id", "name", "sort", "created_at"},
		Fillable:   []string{"name", "slug", "description", "sort", "parent_id"},
		BeforeSave: slugEnsurer(cmsmodels.ArticleCategory{}, "name"),
	}
}

/* ---------------- tags ---------------- */

func tagSpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "tags",
		Model:      cmsmodels.Tag{},
		Searchable: []string{"name", "slug"},
		Sortable:   []string{"id", "name", "created_at"},
		Fillable:   []string{"name", "slug"},
		BeforeSave: slugEnsurer(cmsmodels.Tag{}, "name"),
	}
}

/* ---------------- articles ---------------- */

const maxSlugLen = 191

func articleSpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "articles",
		Model:      cmsmodels.Article{},
		Searchable: []string{"title", "summary"}, // CJK fallback columns
		Sortable:   []string{"id", "title", "status", "published_at", "view_count", "created_at"},
		Fillable: []string{
			"title", "slug", "content", "summary", "cover",
			"category_id", "author_id", "tags", // status is NOT fillable:
			// it changes only through the lifecycle actions.
		},
		Transform:      renderArticle,
		SearchOverride: searchArticles,
		BeforeSave:     prepareArticle,
		AfterSave:      syncTags,

		SoftDeletable: true,
	}
}

// searchArticles: full-text over the generated tsvector for non-CJK terms,
// ILIKE degradation for CJK (plan §M2). The term is always a bound parameter.
func searchArticles(query orm.Query, term string) orm.Query {
	if cmsservices.ContainsCJK(term) {
		like := "%" + strings.ToLower(term) + "%"
		return query.Where(
			"(lower(title) LIKE ? OR lower(coalesce(summary,'')) LIKE ? OR lower(coalesce(content,'')) LIKE ?)",
			like, like, like)
	}
	return query.Where("search_vector @@ websearch_to_tsquery('simple', ?)", term)
}

// prepareArticle validates the cover URL policy and derives a unique slug
// when the author left it blank.
func prepareArticle(data map[string]any) error {
	if cover, present := data["cover"]; present {
		coverURL, _ := cover.(string)
		if err := cmsservices.ValidateCover(coverURL); err != nil {
			return err
		}
	}
	return slugEnsurer(cmsmodels.Article{}, "title")(data)
}

// renderArticle flattens tag links into a newline string so the admin form's
// textarea can round-trip them without a repeater control.
func renderArticle(row any) map[string]any {
	b, err := json.Marshal(row)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{}
	}

	id, ok := toUint64(m["id"])
	if !ok {
		return m
	}
	names, err := articleTagNames(id)
	if err == nil && len(names) > 0 {
		m["tags"] = strings.Join(names, "\n")
	} else {
		m["tags"] = ""
	}
	return m
}

func articleTagNames(articleID uint64) ([]string, error) {
	pivot := cmsmodels.ArticleTag{}
	var ids []uint64
	if err := facades.Orm().Query().
		Table(pivot.TableName()).
		Where("article_id = ?", articleID).
		Pluck("tag_id", &ids); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	values := make([]any, len(ids))
	for i, id := range ids {
		values[i] = id
	}
	var tags []cmsmodels.Tag
	if err := facades.Orm().Query().
		WhereIn("id", values).
		Get(&tags); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(tags))
	for _, t := range tags {
		names = append(names, t.Name)
	}
	return names, nil
}

// syncTags persists the newline-separated tags payload after create/update.
// The gateway hands us a dereferenced model VALUE, hence the value assertion.
func syncTags(instance any, data map[string]any) error {
	article, ok := instance.(cmsmodels.Article)
	if !ok || article.ID == 0 {
		return nil
	}
	raw, present := data["tags"]
	if !present {
		return nil
	}
	text, _ := raw.(string)
	names := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == ','
	})
	return cmsservices.SyncArticleTags(article.ID, names)
}

// slugEnsurer returns a BeforeSave hook that fills a blank slug from the
// resource's human name and guarantees uniqueness against its own table.
func slugEnsurer(model any, nameKey string) func(map[string]any) error {
	return func(data map[string]any) error {
		slug, _ := data["slug"].(string)
		if strings.TrimSpace(slug) != "" {
			return nil
		}
		name, _ := data[nameKey].(string)
		if name == "" {
			return nil
		}
		exists := func(candidate string) (bool, error) {
			var count int64
			count, err := facades.Orm().Query().
				Model(model).
				Where("slug = ?", candidate).
				Count()
			return count > 0, err
		}
		unique, err := cmsservices.EnsureUnique(cmsservices.Slugify(name), maxSlugLen, exists)
		if err != nil {
			return err
		}
		data["slug"] = unique
		return nil
	}
}

func toUint64(v any) (uint64, bool) {
	switch n := v.(type) {
	case float64:
		return uint64(n), true
	case uint64:
		return n, true
	case int64:
		return uint64(n), true
	case json.Number:
		u, err := strconv.ParseUint(n.String(), 10, 64)
		return u, err == nil
	}
	return 0, false
}

/* ---------------- routes ---------------- */

func registerActionRoutes() {
	actions := cmshttp.NewArticleActionController()
	publicAPI := cmshttp.NewPublicController()
	r := facades.Route()

	// Public v1 read endpoints (anonymous, gated by site.mode).
	r.Get("/api/v1/articles", publicAPI.ListArticles)
	r.Get("/api/v1/articles/{slug}", publicAPI.ShowArticle)

	r.Post("/api/admin/articles/{id}/publish", actions.Publish)
	r.Post("/api/admin/articles/{id}/schedule", actions.Schedule)
	r.Post("/api/admin/articles/{id}/unpublish", actions.Unpublish)
	r.Post("/api/admin/articles/{id}/archive", actions.Archive)
	r.Post("/api/admin/articles/{id}/restore", actions.RestoreFromTrash)
	r.Post("/api/admin/articles/{id}/pin", actions.Pin)
	r.Post("/api/admin/articles/{id}/unpin", actions.Unpin)

	// {id}/{action} subtrees shadow the wildcard bulk-delete; re-bind it.
	adminhub.RegisterBulkDeleteFor("articles")
}
