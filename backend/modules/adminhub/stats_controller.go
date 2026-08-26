package adminhub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	authcontrollers "reflexcms/backend/modules/auth/http/controllers"
)

type StatsController struct{}

func NewStatsController() *StatsController { return &StatsController{} }

// Stats: GET /api/admin/stats — dashboard aggregates. Tables belonging to
// later milestones are counted only when they already exist, so the endpoint
// stays green across milestones without feature branches.
func (r *StatsController) Index(ctx http.Context) http.Response {
	if _, ok := authcontrollers.RequireIdentity(ctx); !ok {
		return httpx.Error(ctx, 401, "Unauthorized")
	}

	// Aggregates scan every content table; a short TTL keeps dashboards snappy
	// without stale numbers mattering (plan §M5 performance baseline). The
	// payload travels as a JSON string because redis round-trips lose the
	// concrete map type.
	if raw := facades.Cache().Get(statsCacheKey); raw != nil {
		var cached http.Json
		switch v := raw.(type) {
		case string: // redis round-trip
			if v != "" {
				cached = decodeJSONMap(v)
			}
		default: // memory driver keeps the map
			if b, err := json.Marshal(v); err == nil {
				cached = decodeJSONMap(string(b))
			}
		}
		if len(cached) > 0 {
			return ctx.Response().Success().Json(cached)
		}
	}

	payload := r.aggregate()
	if b, err := json.Marshal(payload); err == nil {
		_ = facades.Cache().Put(statsCacheKey, string(b), statsTTL)
	}

	return ctx.Response().Success().Json(payload)
}

const (
	statsCacheKey = "adminhub:stats"
	statsTTL      = 30 * time.Second
)

func (r *StatsController) aggregate() http.Json {
	return http.Json{
		"usersTotal":        countIfTable("users", nil),
		"usersActive":       countIfTable("users", map[string]any{"status": "active"}),
		"articlesTotal":     countIfTable("articles", nil),
		"articlesPublished": countIfTable("articles", map[string]any{"status": "published"}),
		"topicsTotal":       countIfTable("topics", nil),
		"topicsToday":       topicsToday(),
		"repliesTotal":      countIfTable("replies", nil),
		"pendingComments":   countIfTable("comments", map[string]any{"status": "pending"}),
		"recentArticles":    recentArticles(),
		"publishSeries":     publishSeries(),
	}
}

func decodeJSONMap(raw string) http.Json {
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil || len(m) == 0 {
		return nil
	}
	return http.Json(m)
}

// topicsToday counts topics created since local midnight.
func topicsToday() int64 {
	if !facades.Schema().HasTable("topics") {
		return 0
	}
	midnight := time.Now().Truncate(24 * time.Hour)
	count, err := facadesQuery().Table("topics").
		Where("created_at >= ?", midnight).
		Count()
	if err != nil {
		return 0
	}
	return count
}

func countIfTable(table string, filters map[string]any) int64 {
	if !facades.Schema().HasTable(table) {
		return 0
	}
	query := facadesQuery().Table(table)
	for col, val := range filters {
		query = query.Where(col+" = ?", val)
	}
	count, err := query.Count()
	if err != nil {
		return 0
	}
	return count
}

// recentArticles lists the five latest articles for the dashboard widget.
func recentArticles() []map[string]any {
	rows := []map[string]any{}
	if !facades.Schema().HasTable("articles") {
		return rows
	}

	err := facadesQuery().Table("articles").
		Select("id", "title", "slug", "status", "created_at").
		OrderBy("created_at", "desc").
		Limit(5).
		Get(&rows)
	if err != nil {
		return []map[string]any{}
	}
	return rows
}

// publishSeries aggregates published articles per week over the last ten
// weeks, oldest first — the shape consumed by the trend widget.
func publishSeries() []map[string]any {
	series := make([]map[string]any, 0, 10)
	if !facades.Schema().HasTable("articles") {
		for w := 9; w >= 0; w-- {
			series = append(series, map[string]any{
				"label": fmt.Sprintf("W-%d", w), "count": 0,
			})
		}
		return series
	}

	now := time.Now()
	for w := 9; w >= 0; w-- {
		start := now.AddDate(0, 0, -(w+1)*7)
		end := now.AddDate(0, 0, -w*7)

		var count int64
		count, err := facadesQuery().Table("articles").
			Where("published_at >= ? AND published_at < ?", start, end).
			Count()
		if err != nil {
			count = 0
		}
		series = append(series, map[string]any{
			"label": fmt.Sprintf("W-%d", w),
			"start": start.Format("2006-01-02"),
			"count": count,
		})
	}
	return series
}
