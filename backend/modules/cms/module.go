package cms

import (
	"time"

	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/schedule"

	"reflexcms/backend/app/facades"
	"reflexcms/backend/database/migrations"
	cmsmodels "reflexcms/backend/modules/cms/models"
	cmsservices "reflexcms/backend/modules/cms/services"
)

// Module is the CMS feature module (articles, categories, tags).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "cms" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.CmsMigrations()
}

func (m *Module) Routes() {
	Register()
}

// Schedule publishes due scheduled articles every minute (WordPress-style
// timed releases). Products without cms never schedule it.
func (m *Module) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Call(func() {
			now := time.Now()
			var due []cmsmodels.Article
			if err := facades.Orm().Query().
				Where("status = ? AND published_at IS NOT NULL AND published_at <= ?",
					cmsmodels.ArticleScheduled, now).
				Get(&due); err != nil {
				facades.Log().Warning("cms: scheduled publish query failed: " + err.Error())
				return
			}
			for i := range due {
				if err := cmsservices.Transition(&due[i], cmsmodels.ArticlePublished); err != nil {
					continue
				}
				facades.Log().Info("cms: scheduled article #" + itoa(int(due[i].ID)) + " published")
			}
		}).EveryMinute(),
	}
}

func itoa(n int) string {
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	if digits == "" {
		return "0"
	}
	return digits
}

func (m *Module) Boot() {}
