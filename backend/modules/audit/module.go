package audit

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/database/migrations"
	adminhub "reflexcms/backend/modules/adminhub"
	auditmodels "reflexcms/backend/modules/audit/models"
)

// Module exposes the admin operation audit trail (read-only).
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Name() string { return "audit" }

func (m *Module) Migrations() []schema.Migration {
	return migrations.AuditMigrations()
}

func (m *Module) Routes() {
	adminhub.Register(&adminhub.Spec{
		Name:             "operation_logs",
		Model:            auditmodels.AdminOperationLog{},
		Searchable:       []string{"action", "resource"},
		Sortable:         []string{"id", "user_id", "action", "resource", "created_at"},
		Fillable:         []string{}, // system-written: no direct create/update
		PermissionPrefix: "audit",
	})
}

func (m *Module) Boot() {}
