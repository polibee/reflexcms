package adminhub

import (
	"github.com/goravel/framework/contracts/http"

	authcontrollers "reflexcms/backend/modules/auth/http/controllers"
)

// CheckPermission resolves the caller and enforces the required permission.
// It returns nil when the call may proceed, otherwise a ready-to-return
// 401/403 response. Module controllers reuse this for their action routes.
func CheckPermission(ctx http.Context, required string) *http.Response {
	return requirePermission(ctx, required)
}

// ParseRowID extracts and validates the {id} route segment.
func ParseRowID(ctx http.Context) (uint64, bool) {
	return parseID(ctx)
}

// CurrentUserID resolves the caller's id (0 when anonymous).
func CurrentUserID(ctx http.Context) uint64 {
	if identity, ok := authcontrollers.RequireIdentity(ctx); ok {
		return identity.ID
	}
	return 0
}
