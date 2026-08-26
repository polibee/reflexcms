package http

import (
	"github.com/goravel/framework/contracts/http"

	httpx "reflexcms/backend/app/support/httpx"
	authcontrollers "reflexcms/backend/modules/auth/http/controllers"
	forumservices "reflexcms/backend/modules/forum/services"
)

type LikeController struct{}

func NewLikeController() *LikeController { return &LikeController{} }

// Toggle: POST /api/admin/forum/like — {likeable_type, likeable_id}.
// Any authenticated admin user may vote; the service keeps counters exact.
func (r *LikeController) Toggle(ctx http.Context) http.Response {
	if _, ok := authcontrollers.RequireIdentity(ctx); !ok {
		return httpx.Error(ctx, 401, "Unauthorized")
	}

	var req struct {
		LikeableType string `json:"likeable_type"`
		LikeableID   uint64 `json:"likeable_id"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.LikeableID == 0 || req.LikeableType == "" {
		return httpx.Error(ctx, 422, "likeable_type and likeable_id are required")
	}

	liked, count, err := forumservices.ToggleLike(identityID(ctx), req.LikeableType, req.LikeableID)
	if err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}
	return ctx.Response().Success().Json(http.Json{
		"liked": liked, "like_count": count,
	})
}

func identityID(ctx http.Context) uint64 {
	identity, ok := authcontrollers.RequireIdentity(ctx)
	if !ok {
		return 0
	}
	return identity.ID
}
