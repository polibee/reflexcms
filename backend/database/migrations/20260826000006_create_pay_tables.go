package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"reflexcms/backend/app/facades"
)

type M20260826000006CreateProductsTable struct{}

func (r *M20260826000006CreateProductsTable) Signature() string {
	return "20260826000006_create_products_table"
}

func (r *M20260826000006CreateProductsTable) Up() error {
	if facades.Schema().HasTable("products") {
		return nil
	}
	return facades.Schema().Create("products", func(table schema.Blueprint) {
		table.ID()
		table.String("title", 255)
		table.Text("description").Nullable()
		table.Integer("price_cents").Default(0)
		table.String("currency", 8).Default("USD")
		table.Integer("stock").Default(0)
		table.Integer("grant_points").Default(0)
		table.String("image", 512).Nullable()
		table.Boolean("is_active").Default(true)
		table.Integer("sort").Default(100)
		table.TimestampsTz()
		table.Index("is_active")
	})
}

func (r *M20260826000006CreateProductsTable) Down() error {
	return facades.Schema().Drop("products")
}

type M20260826000007CreateOrdersTable struct{}

func (r *M20260826000007CreateOrdersTable) Signature() string {
	return "20260826000007_create_orders_table"
}

func (r *M20260826000007CreateOrdersTable) Up() error {
	if facades.Schema().HasTable("orders") {
		return nil
	}
	return facades.Schema().Create("orders", func(table schema.Blueprint) {
		table.ID()
		table.String("order_no", 64)
		table.UnsignedBigInteger("user_id")
		table.UnsignedBigInteger("product_id").Nullable()
		table.String("title", 255)
		table.Integer("amount_cents").Default(0)
		table.String("currency", 8).Default("USD")
		table.String("gateway", 32).Default("")
		table.String("gateway_ref", 128).Nullable()
		table.String("status", 32).Default("pending")
		table.DateTimeTz("paid_at").Nullable()
		table.TimestampsTz()
		table.Unique("order_no")
		table.Index("user_id")
		table.Index("status")
	})
}

func (r *M20260826000007CreateOrdersTable) Down() error {
	return facades.Schema().Drop("orders")
}

type M20260826000008CreateBoardModeratorsTable struct{}

func (r *M20260826000008CreateBoardModeratorsTable) Signature() string {
	return "20260826000008_create_board_moderators_table"
}

func (r *M20260826000008CreateBoardModeratorsTable) Up() error {
	if facades.Schema().HasTable("board_moderators") {
		return nil
	}
	return facades.Schema().Create("board_moderators", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("board_id")
		table.UnsignedBigInteger("user_id")
		table.TimestampsTz()
		table.Unique("board_id", "user_id")
		table.Index("user_id")
	})
}

func (r *M20260826000008CreateBoardModeratorsTable) Down() error {
	return facades.Schema().Drop("board_moderators")
}
