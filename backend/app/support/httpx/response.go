package httpx

import (
	"github.com/goravel/framework/contracts/http"
)

// Error writes the failure envelope expected by the nuxtadmin BFF: h3 picks
// statusCode/statusMessage out of the JSON body when rethrowing on the client,
// so notifyError surfaces statusMessage verbatim (plan §5.2 error semantics).
func Error(ctx http.Context, code int, message string) http.Response {
	return ctx.Response().Json(code, http.Json{
		"statusCode":    code,
		"statusMessage": message,
	})
}
