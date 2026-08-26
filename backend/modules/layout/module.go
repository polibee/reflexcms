package layout

import (
	"encoding/json"
	"strings"

	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	"reflexcms/backend/database/migrations"
	adminhub "reflexcms/backend/modules/adminhub"
	layoutmodels "reflexcms/backend/modules/layout/models"
)

// Module is the Layout feature module: sidebar widgets, menus, carousels,
// and invite codes — all manageable from the admin panel.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "layout" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.LayoutMigrations()
}

func (m *Module) Routes() {
	registerAdminSpecs()
	registerPublicRoutes()
}

func (m *Module) Boot() {}

// widgetConfigForForm unwraps custom_html's {"html": ...} storage back into
// the raw HTML string so the admin form edits code, not JSON wrapping.
func widgetConfigForForm(row any) map[string]any {
	w, ok := row.(layoutmodels.Widget)
	if !ok {
		return map[string]any{}
	}
	data, err := json.Marshal(w)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]any{}
	}
	if w.Type == "custom_html" {
		var cfg map[string]any
		if json.Unmarshal([]byte(w.Config), &cfg) == nil {
			if html, ok := cfg["html"].(string); ok {
				m["config"] = html
			}
		}
	}
	return m
}

// presetConfigDefaults fill in when a widget is saved with an empty config,
// so picking a preset card type "just works" without touching JSON.
var presetConfigDefaults = map[string]string{
	"categories":       `{"show_count": true, "count": 10}`,
	"tags":             `{"max_tags": 20}`,
	"latest_articles":  `{"count": 5}`,
	"hottest_articles": `{"count": 5}`,
	"latest_topics":    `{"count": 5}`,
	"hottest_topics":   `{"count": 5}`,
	"boards":           `{}`,
	"user_stats":       `{}`,
	"author_center":    `{}`,
}

// normalizeWidgetConfig stores raw custom_html input as {"html": ...}; JSON
// payloads (hand-written or previously stored) are kept verbatim. The wrap
// is deliberately type-agnostic: partial updates may omit widget_type, and
// the config column is jsonb — a non-JSON payload could not save any other
// way, so the only sensible reading is "this is raw custom HTML". Preset
// types with an empty config get their working defaults.
func normalizeWidgetConfig(data map[string]any) error {
	raw, ok := data["config"]
	if !ok {
		return nil
	}
	cfg, ok := raw.(string)
	if !ok {
		return nil
	}
	cfg = strings.TrimSpace(cfg)

	if typ, _ := data["widget_type"].(string); cfg == "" {
		if def, known := presetConfigDefaults[typ]; known {
			data["config"] = def
		}
		return nil
	}

	if json.Valid([]byte(cfg)) {
		return nil
	}
	encoded, err := json.Marshal(map[string]string{"html": cfg})
	if err != nil {
		return err
	}
	data["config"] = string(encoded)
	return nil
}

func registerAdminSpecs() {
	adminhub.Register(&adminhub.Spec{
		Name:       "widgets",
		Model:      layoutmodels.Widget{},
		Searchable: []string{"title", "widget_type"},
		Sortable:   []string{"id", "area", "sort", "created_at"},
		Fillable:   []string{"area", "widget_type", "title", "config", "sort", "is_active"},
		// custom_html accepts raw HTML/JS in the config box: anything that is
		// not valid JSON is wrapped as {"html": ...} on save and unwrapped
		// again for editing. Explicit JSON still passes through untouched.
		Transform:  widgetConfigForForm,
		BeforeSave: normalizeWidgetConfig,
	})

	adminhub.Register(&adminhub.Spec{
		Name:       "menus",
		Model:      layoutmodels.Menu{},
		Searchable: []string{"label", "url"},
		Sortable:   []string{"id", "location", "sort", "created_at"},
		Fillable:   []string{"location", "label", "url", "sort", "parent_id", "is_visible"},
	})

	adminhub.Register(&adminhub.Spec{
		Name:       "carousels",
		Model:      layoutmodels.Carousel{},
		Searchable: []string{"title"},
		Sortable:   []string{"id", "sort", "created_at"},
		Fillable:   []string{"title", "image_url", "link_url", "sort", "is_active"},
	})
}

// registerPublicRoutes wires anonymous /api/v1 endpoints consumed by the
// public frontend for layout rendering.
func registerPublicRoutes() {
	r := facades.Route()

	// Sidebar widgets grouped by area.
	r.Get("/api/v1/widgets", func(ctx http.Context) http.Response {
		area := ctx.Request().Query("area", "")
		if area == "" {
			return httpx.Error(ctx, 422, "area parameter is required")
		}
		items := []map[string]any{}
		if err := facades.Orm().Query().Table("widgets").
			Where("area = ? AND is_active = ?", area, true).
			Select("id", "area", "widget_type", "title", "config", "sort").
			OrderBy("sort", "asc").
			Get(&items); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"items": items})
	})

	// Menu tree per location.
	r.Get("/api/v1/menus", func(ctx http.Context) http.Response {
		location := ctx.Request().Query("location", "")
		if location == "" {
			return httpx.Error(ctx, 422, "location parameter is required")
		}
		items := []map[string]any{}
		if err := facades.Orm().Query().Table("menus").
			Where("location = ? AND is_visible = ?", location, true).
			Select("id", "label", "url", "sort", "parent_id").
			OrderBy("sort", "asc").
			Get(&items); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		return ctx.Response().Success().Json(http.Json{"items": items})
	})

	// Active carousel slides.
	r.Get("/api/v1/carousels", func(ctx http.Context) http.Response {
		var items []map[string]any
		if err := facades.Orm().Query().Table("carousels").
			Where("is_active = ?", true).
			Select("id", "title", "image_url", "link_url", "sort").
			OrderBy("sort", "asc").
			Get(&items); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		if items == nil {
			items = []map[string]any{}
		}
		return ctx.Response().Success().Json(http.Json{"items": items})
	})

	// Sidebar data: returns each active widget with its rendered data
	// (categories list, latest articles, user stats, etc.) so the frontend
	// makes ONE call to render the full sidebar.
	r.Get("/api/v1/sidebar", func(ctx http.Context) http.Response {
		area := ctx.Request().Query("area", "")
		if area == "" {
			return httpx.Error(ctx, 422, "area parameter is required")
		}

		var widgets []map[string]any
		if err := facades.Orm().Query().Table("widgets").
			Where("area = ? AND is_active = ?", area, true).
			Select("id", "widget_type", "title", "config", "sort").
			OrderBy("sort", "asc").
			Get(&widgets); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}

		result := []map[string]any{}
		for _, w := range widgets {
			widgetType, _ := w["widget_type"].(string)
			entry := map[string]any{"type": widgetType, "title": w["title"]}

			switch widgetType {
			case "categories":
				var cats []map[string]any
				facades.Orm().Query().Table("article_categories").
					Select("id", "name", "slug").
					OrderBy("sort", "asc").Get(&cats)
				entry["data"] = cats
			case "tags":
				var tags []map[string]any
				facades.Orm().Query().Table("tags").
					Select("id", "name", "slug").
					OrderBy("name", "asc").Limit(30).Get(&tags)
				entry["data"] = tags
			case "latest_articles":
				var arts []map[string]any
				facades.Orm().Query().Table("articles").
					Where("status = 'published' AND deleted_at IS NULL").
					Select("id", "title", "slug", "cover", "published_at").
					OrderBy("published_at", "desc").Limit(5).Get(&arts)
				entry["data"] = arts
			case "hottest_articles":
				var arts []map[string]any
				facades.Orm().Query().Table("articles").
					Where("status = 'published' AND deleted_at IS NULL").
					Select("id", "title", "slug", "cover", "view_count").
					OrderBy("view_count", "desc").Limit(5).Get(&arts)
				entry["data"] = arts
			case "latest_topics":
				var tps []map[string]any
				facades.Orm().Query().Table("topics").
					Where("status = 'open' AND deleted_at IS NULL").
					Select("id", "title", "reply_count").
					OrderBy("created_at", "desc").Limit(5).Get(&tps)
				entry["data"] = tps
			case "hottest_topics":
				var tps []map[string]any
				facades.Orm().Query().Table("topics").
					Where("status = 'open' AND deleted_at IS NULL").
					Select("id", "title", "view_count").
					OrderBy("view_count", "desc").Limit(5).Get(&tps)
				entry["data"] = tps
			case "boards":
				var boards []map[string]any
				facades.Orm().Query().Table("forum_categories").
					Select("id", "name", "slug", "topic_count").
					OrderBy("sort", "asc").Get(&boards)
				entry["data"] = boards
			case "user_stats":
				userCount, _ := facades.Orm().Query().Table("users").Count()
				topicCount, _ := facades.Orm().Query().Table("topics").Count()
				entry["data"] = map[string]any{"users": userCount, "topics": topicCount}
			case "author_center":
				// Identity-dependent; the card body renders client-side.
				entry["data"] = map[string]any{}
			case "composer":
				// Generic post-entry card; renders a prominent 发帖 button.
				entry["data"] = map[string]any{}
			case "custom_html", "custom_image", "custom_text_link":
				var cfg map[string]any
				if cfgStr, ok := w["config"].(string); ok && cfgStr != "" {
					json.Unmarshal([]byte(cfgStr), &cfg)
				}
				entry["data"] = cfg
			}

			result = append(result, entry)
		}

		// gorm leaves target slices nil when a query matches no rows; marshal
		// those as [] so public consumers never see null lists.
		for _, e := range result {
			if list, ok := e["data"].([]map[string]any); ok && list == nil {
				e["data"] = []map[string]any{}
			}
		}

		return ctx.Response().Success().Json(http.Json{"widgets": result})
	})
}
