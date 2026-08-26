package access

import (
	"github.com/goravel/framework/contracts/hash"

	"reflexcms/backend/app/facades"
)

func facadesHash() hash.Hash {
	return facades.Hash()
}
