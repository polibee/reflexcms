package controllers

import (
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"

	httpx "reflexcms/backend/app/support/httpx"
	authservices "reflexcms/backend/modules/auth/services"
)

type AuthController struct{}

func NewAuthController() *AuthController { return &AuthController{} }

// Login: POST /api/auth/login — {email,password} → {token,user}.
// The admin BFF stores the token in its own httpOnly cookie.
func (r *AuthController) Login(ctx http.Context) http.Response {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.Request().Bind(&req); err != nil || req.Email == "" || req.Password == "" {
		return httpx.Error(ctx, 422, "Email and password are required")
	}

	token, identity, err := authservices.Login(
		req.Email, req.Password,
		ctx.Request().Ip(), ctx.Request().Header("User-Agent", ""),
	)
	switch {
	case errors.Is(err, authservices.ErrAccountLocked),
		errors.Is(err, authservices.ErrRateLimited):
		return httpx.Error(ctx, http.StatusTooManyRequests, err.Error())
	case errors.Is(err, authservices.ErrInvalidCredentials):
		return httpx.Error(ctx, http.StatusUnauthorized, "Invalid credentials")
	case err != nil:
		return httpx.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Response().Success().Json(http.Json{
		"token": token,
		"user":  identity,
	})
}

// Me: GET /api/auth/me → current AuthUser (or 401).
func (r *AuthController) Me(ctx http.Context) http.Response {
	identity, ok := RequireIdentity(ctx)
	if !ok {
		return httpx.Error(ctx, http.StatusUnauthorized, string(authservices.ErrSessionExpired.Error()))
	}
	return ctx.Response().Success().Json(identity)
}

// Logout: POST /api/auth/logout — revokes the presented session.
func (r *AuthController) Logout(ctx http.Context) http.Response {
	if token := BearerToken(ctx); token != "" {
		_ = authservices.Revoke(token)
	}
	return ctx.Response().Success().Json(http.Json{"ok": true})
}

/* ---------------- shared helpers ---------------- */

// BearerToken extracts the opaque token from the Authorization header.
func BearerToken(ctx http.Context) string {
	header := ctx.Request().Header("Authorization", "")
	const prefix = "Bearer "
	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}
	return ""
}

// RequireIdentity resolves the caller from the bearer token. ok=false means
// unauthenticated; handlers respond with 401 themselves so every endpoint
// keeps its error envelope explicit.
func RequireIdentity(ctx http.Context) (*authservices.Identity, bool) {
	token := BearerToken(ctx)
	if token == "" {
		return nil, false
	}
	identity, err := authservices.Resolve(token)
	if err != nil {
		return nil, false
	}
	return identity, true
}
