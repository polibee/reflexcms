package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260826000005CreateBlockedUsersTable struct{}

func (r *M20260826000005CreateBlockedUsersTable) Signature() string {
	return "20260826000005_create_blocked_users_table"
}

func (r *M20260826000005CreateBlockedUsersTable) Up() error {
	if facades.Schema().HasTable("blocked_users") {
		return nil
	}
	return facades.Schema().Create("blocked_users", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.UnsignedBigInteger("blocked_id")
		table.DateTimeTz("created_at").Nullable()
		table.Unique("user_id", "blocked_id")
		table.Index("user_id")
	})
}

func (r *M20260826000005CreateBlockedUsersTable) Down() error {
	return facades.Schema().Drop("blocked_users")
}
