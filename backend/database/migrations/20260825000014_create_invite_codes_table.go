package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000014CreateInviteCodesTable struct{}

func (r *M20260825000014CreateInviteCodesTable) Signature() string {
	return "20260825000014_create_invite_codes_table"
}

func (r *M20260825000014CreateInviteCodesTable) Up() error {
	if facades.Schema().HasTable("invite_codes") {
		return nil
	}
	return facades.Schema().Create("invite_codes", func(table schema.Blueprint) {
		table.ID()
		table.String("code", 32)
		table.UnsignedBigInteger("creator_id")
		table.Integer("max_uses").Default(1)
		table.Integer("used_count").Default(0)
		table.DateTimeTz("expires_at").Nullable()
		table.TimestampsTz()
		table.Unique("code")
		table.Index("creator_id")
	})
}

func (r *M20260825000014CreateInviteCodesTable) Down() error {
	return facades.Schema().DropIfExists("invite_codes")
}
