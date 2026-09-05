package access

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	"reflexcms/backend/modules/access/models"
	accessservices "reflexcms/backend/modules/access/services"
	adminhub "reflexcms/backend/modules/adminhub"
)

// Register exposes the RBAC domain through the admin gateway.
func Register() {
	adminhub.Register(&adminhub.Spec{
		Name:       "users",
		Model:      models.User{},
		Searchable: []string{"username", "email"},
		Sortable:   []string{"id", "username", "email", "status", "created_at"},
		Fillable:   []string{"username", "email", "password", "avatar", "bio", "website", "role_id", "status"},
		Hidden:     []string{"password"},
		BeforeSave: hashPassword,

		SoftDeletable: true,
	})

	adminhub.Register(&adminhub.Spec{
		Name:             "roles",
		Model:            models.Role{},
		Searchable:       []string{"name", "display_name"},
		Sortable:         []string{"id", "name", "sort", "created_at"},
		Fillable:         []string{"name", "display_name", "permissions", "sort"},
		Transform:        renderRole,
		BeforeSave:       prepareRole,
		PermissionPrefix: "roles",
	})

	registerRoutes()
}

func registerRoutes() {
	// Permission catalog powering the visual role editor (any signed-in admin).
	facades.Route().Get("/api/admin/permissions", func(ctx http.Context) http.Response {
		return ctx.Response().Success().Json(accessservices.Catalog())
	})

	// Public user homepages + self-service profile/security + favorites.
	r := facades.Route()
	r.Get("/api/v1/users/{id}", ShowUser)
	r.Get("/api/v1/users/{id}/topics", UserTopics)
	r.Get("/api/v1/users/{id}/replies", UserReplies)
	r.Get("/api/v1/users/{id}/favorites", UserFavorites)
	r.Get("/api/v1/me/profile", MyProfile)
	r.Get("/api/v1/me/stats", MyStats)
	r.Get("/api/v1/me/points", MyPoints)
	r.Post("/api/v1/me/signin", SignIn)
	r.Put("/api/v1/me/profile", UpdateMyProfile)
	r.Put("/api/v1/me/profile-card", UpdateMyProfileCard)
	r.Put("/api/v1/me/password", ChangeMyPassword)
	r.Post("/api/v1/favorites/toggle", ToggleFavorite)
	r.Get("/api/v1/favorites/state", MyFavoriteState)
	r.Post("/api/v1/me/messages", SendMessage)
	r.Get("/api/v1/me/messages", MyMessages)
	r.Get("/api/v1/me/messages/unread-count", UnreadMessages)
	r.Post("/api/v1/me/messages/{id}/read", MarkMessageRead)
	r.Get("/api/v1/me/notifications", MyNotifications)
	r.Post("/api/v1/me/notifications/read-all", ReadAllNotifications)
	r.Post("/api/v1/me/blocked/toggle", ToggleBlock)
	r.Get("/api/v1/me/blocked", MyBlocked)
	r.Post("/api/v1/auth/register", RegisterPublic)
	r.Get("/api/v1/auth/register-config", func(ctx http.Context) http.Response {
		return ctx.Response().Success().Json(http.Json{
			"invite_required": regInviteRequired(),
			"enabled":         regEnabled(),
		})
	})
}

// renderRole passes permissions through as a JSON string[] so the
// checkboxGroup field on the role form binds and prefills directly.
func renderRole(row any) map[string]any {
	b, err := json.Marshal(row)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{}
	}
	return m
}

// prepareRole validates every concrete permission against the catalog and
// normalizes the incoming shape: the visual checkboxGroup submits string[],
// the legacy textarea submits newline-separated text.
func prepareRole(data map[string]any) error {
	raw, present := data["permissions"]
	if !present {
		return nil
	}

	var entries []string
	switch v := raw.(type) {
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				entries = append(entries, strings.TrimSpace(s))
			}
		}
	case string:
		entries = strings.FieldsFunc(v, func(r rune) bool { return r == '\n' || r == ',' || r == ' ' })
	}

	out := make([]string, 0, len(entries))
	for _, f := range entries {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !accessservices.IsKnownPermission(f) {
			return fmt.Errorf("unknown permission: %s", f)
		}
		out = append(out, f)
	}
	data["permissions"] = out
	return nil
}

// hashPassword turns a plaintext password payload into its bcrypt hash before
// the value ever reaches the model layer. Empty passwords are dropped so an
// edit that leaves the field untouched does not clobber the stored hash.
func hashPassword(data map[string]any) error {
	raw, present := data["password"]
	if !present {
		return nil
	}
	plain, _ := raw.(string)
	if plain == "" {
		delete(data, "password")
		return nil
	}
	hashed, err := facadesHash().Make(plain)
	if err != nil {
		return err
	}
	data["password"] = hashed
	return nil
}
