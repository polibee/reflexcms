package settings

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
	"reflexcms/backend/database/migrations"
	adminhub "reflexcms/backend/modules/adminhub"
	settingshttp "reflexcms/backend/modules/settings/http"
)

// Module is the site configuration feature module (KV store, site.mode,
// public /api/v1/site/config).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "settings" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.SettingsMigrations()
}

func (m *Module) Routes() {
	adminhub.Register(settingSpec())

	facades.Route().Get("/api/v1/site/config", settingshttp.NewSiteConfigController().Show)
}

func (m *Module) Boot() {
	seedDefaults()
}
