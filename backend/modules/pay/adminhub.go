package pay

import (
	"github.com/goravel/framework/contracts/http"

	adminhub "reflexcms/backend/modules/adminhub"
)

// Adapters over the adminhub gateway so the rest of the module stays terse.

type adminhubSpec = adminhub.Spec

func adminhubRegister(spec *adminhubSpec) { adminhub.Register(spec) }

func adminhubCheckPermission(ctx http.Context, permission string) *http.Response {
	if resp := adminhub.CheckPermission(ctx, permission); resp != nil {
		return resp
	}
	return nil
}
