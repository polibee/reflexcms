package main

import (
	"reflexcms/backend/bootstrap"
	"reflexcms/backend/products"
)

// forum is the boards-only product: platform core + Forum (boards, topics,
// replies, likes) + notifications + settings. No CMS tables exist here.
func main() {
	app := bootstrap.Boot(products.Forum()...)

	app.Start()
}
