package forum

import (
	"strings"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	"reflexcms/backend/app/support/slug"
	adminhub "reflexcms/backend/modules/adminhub"
	forumhttp "reflexcms/backend/modules/forum/http"
	forummodels "reflexcms/backend/modules/forum/models"
	forumservices "reflexcms/backend/modules/forum/services"
)

const maxSlugLen = 191

// Register exposes the Forum domain: three gateway resources, moderation
// action routes and the transactional reply endpoint.
func Register() {
	adminhub.Register(forumCategorySpec())
	adminhub.Register(topicSpec())
	adminhub.Register(replySpec())
	adminhub.Register(boardModeratorSpec())
	registerRoutes()
}

// boardModeratorSpec exposes per-board moderator assignment to the admin.
func boardModeratorSpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:             "board_moderators",
		Model:            forummodels.BoardModerator{},
		Sortable:         []string{"id", "board_id", "user_id", "created_at"},
		Fillable:         []string{"board_id", "user_id"},
		PermissionPrefix: "board_moderators",
	}
}

func forumCategorySpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "forums",
		Model:      forummodels.ForumCategory{},
		Searchable: []string{"name", "slug"},
		Sortable:   []string{"id", "name", "sort", "topic_count", "created_at"},
		Fillable:   []string{"name", "slug", "description", "icon", "sort"},
		BeforeSave: slugEnsurer(forummodels.ForumCategory{}, "name"),
	}
}

func topicSpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "topics",
		Model:      forummodels.Topic{},
		Searchable: []string{"title"},
		Sortable: []string{
			"id", "title", "status", "reply_count", "like_count",
			"view_count", "last_reply_at", "created_at",
		},
		SortRaw: map[string]string{
			// Hot score; weights are developer constants (plan §M3).
			"hot": "view_count*1 + reply_count*3 + like_count*5 + CASE WHEN is_pinned THEN 1000 ELSE 0 END",
		},
		Fillable: []string{"user_id", "forum_category_id", "title", "content"},
		OnShow:   func(instance any) { _ = forumservices.IncrView(instance.(forummodels.Topic).ID) },

		SoftDeletable: true,
	}
}

func replySpec() *adminhub.Spec {
	return &adminhub.Spec{
		Name:       "replies",
		Model:      forummodels.Reply{},
		Searchable: []string{"content"},
		Sortable:   []string{"id", "floor", "like_count", "created_at"},
		Fillable:   []string{"topic_id", "user_id", "parent_id", "content"},

		SoftDeletable: true,
	}
}

/* ---------------- routes ---------------- */

type replyRequest struct {
	TopicID  uint64 `json:"topic_id"`
	UserID   uint64 `json:"user_id"`
	ParentID uint64 `json:"parent_id"`
	Content  string `json:"content"`
}

// registerRoutes wires moderation actions and the transactional reply
// endpoint. Replies created here go through ReplyToTopic so floors and
// counters stay consistent; plain gateway CRUD on /api/admin/replies remains
// available for administrative edits.
func registerRoutes() {
	actions := forumhttp.NewTopicActionController()
	bestReply := forumhttp.NewBestReplyController()
	likes := forumhttp.NewLikeController()
	publicAPI := forumhttp.NewPublicController()
	r := facades.Route()

	// Public v1 read endpoints (anonymous, gated by site.mode) plus the
	// session-authenticated reply endpoint for the public frontend.
	r.Get("/api/v1/forums", publicAPI.ListForums)
	r.Get("/api/v1/topics", publicAPI.ListTopics)
	r.Post("/api/v1/topics", publicAPI.StoreTopic)
	r.Get("/api/v1/topics/{id}", publicAPI.ShowTopic)
	r.Post("/api/v1/topics/{id}/moderate", publicAPI.ModerateTopic)
	r.Get("/api/v1/topics/{id}/replies", publicAPI.TopicReplies)
	r.Get("/api/v1/replies/latest", publicAPI.LatestReplies)
	r.Post("/api/v1/topics/{id}/replies", publicAPI.StoreReply)

	r.Post("/api/admin/topics/{id}/pin", actions.Pin)
	r.Post("/api/admin/topics/{id}/unpin", actions.Unpin)
	r.Post("/api/admin/topics/{id}/feature", actions.Feature)
	r.Post("/api/admin/topics/{id}/unfeature", actions.Unfeature)
	r.Post("/api/admin/topics/{id}/close", actions.Close)
	r.Post("/api/admin/topics/{id}/open", actions.Open)
	r.Post("/api/admin/topics/{id}/best-reply", bestReply.SetBestReply)
	r.Post("/api/admin/forum/like", likes.Toggle)

	// {id}/{action} subtrees shadow the wildcard bulk-delete; re-bind it.
	adminhub.RegisterBulkDeleteFor("topics")

	r.Post("/api/admin/forum/reply", func(ctx http.Context) http.Response {
		if resp := adminhub.CheckPermission(ctx, "replies.create"); resp != nil {
			return *resp
		}

		var req replyRequest
		if err := ctx.Request().Bind(&req); err != nil || strings.TrimSpace(req.Content) == "" || req.TopicID == 0 {
			return httpx.Error(ctx, 422, "topic_id and content are required")
		}

		reply, err := forumservices.ReplyToTopic(req.TopicID, req.UserID, req.ParentID, req.Content)
		if err != nil {
			return httpx.Error(ctx, 422, err.Error())
		}
		return ctx.Response().Success().Json(reply)
	})
}

/* ---------------- shared helpers ---------------- */

func slugEnsurer(model any, nameKey string) func(map[string]any) error {
	return func(data map[string]any) error {
		slugVal, _ := data["slug"].(string)
		if strings.TrimSpace(slugVal) != "" {
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
		unique, err := slug.EnsureUnique(slug.Slugify(name), maxSlugLen, exists)
		if err != nil {
			return err
		}
		data["slug"] = unique
		return nil
	}
}
