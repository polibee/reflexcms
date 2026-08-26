package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260825000015CreatePointsTables struct{}

func (r *M20260825000015CreatePointsTables) Signature() string {
	return "20260825000015_create_points_tables"
}

func (r *M20260825000015CreatePointsTables) Up() error {
	if !facades.Schema().HasTable("user_points") {
		if err := facades.Schema().Create("user_points", func(table schema.Blueprint) {
			table.UnsignedBigInteger("user_id")
			table.Integer("balance").Default(0)
			table.Integer("total_earned").Default(0)
			table.DateTimeTz("updated_at")
			table.Primary("user_id")
		}); err != nil {
			return err
		}
	}

	if facades.Schema().HasTable("point_transactions") {
		return nil
	}
	return facades.Schema().Create("point_transactions", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.Integer("amount")
		table.String("reason", 64)
		table.TimestampsTz()
		table.Index("user_id")
	})
}

func (r *M20260825000015CreatePointsTables) Down() error {
	if err := facades.Schema().DropIfExists("point_transactions"); err != nil {
		return err
	}
	return facades.Schema().DropIfExists("user_points")
}
