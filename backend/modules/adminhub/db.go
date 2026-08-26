package adminhub

import (
	"github.com/goravel/framework/contracts/database/orm"

	"reflexcms/backend/app/facades"
)

// facadesQuery centralises the ORM entry point so every gateway query shares
// one construction site. All external values are bound parameters; identifier
// interpolation only ever uses whitelist-validated names.
func facadesQuery() orm.Query {
	return facades.Orm().Query()
}
