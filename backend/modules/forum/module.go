package forum

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/schedule"

	"reflexcms/backend/app/facades"
	"reflexcms/backend/database/migrations"
	forumservices "reflexcms/backend/modules/forum/services"
)

// Module is the Forum feature module (boards, topics, replies, likes).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "forum" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.ForumMigrations()
}

func (m *Module) Routes() {
	Register()
}

func (m *Module) Boot() {}

// Schedule owns the buffered-view flush: Redis counters are drained into
// topics.view_count every five minutes. Products without forum never
// schedule it.
func (m *Module) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule().Call(func() {
			if err := forumservices.FlushViews(); err != nil {
				facades.Log().Warning("forum: view flush failed: " + err.Error())
			}
		}).EveryFiveMinutes(),
	}
}
