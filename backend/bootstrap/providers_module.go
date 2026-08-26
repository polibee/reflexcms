package bootstrap

import (
	"github.com/goravel/framework/contracts/foundation"

	"reflexcms/backend/kernel"
)

// moduleServiceProvider drives the composed modules through their Boot phase
// in composition order, after all framework providers are up.
type moduleServiceProvider struct {
	mods []kernel.Module
}

func (p *moduleServiceProvider) Register(app foundation.Application) {}

func (p *moduleServiceProvider) Boot(app foundation.Application) {
	for _, m := range p.mods {
		m.Boot()
	}
}

var _ foundation.ServiceProvider = (*moduleServiceProvider)(nil)
