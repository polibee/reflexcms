package comment

import (
	"github.com/goravel/framework/contracts/route"

	"reflexcms/backend/app/facades"
)

func facadesRoute() route.Route {
	return facades.Route()
}
