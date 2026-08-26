package main

import (
	"reflexcms/backend/bootstrap"
	"reflexcms/backend/products"
)

// blog is the articles-only product: platform core + CMS + comments +
// notifications + settings. No forum tables, routes or schedules exist here.
func main() {
	app := bootstrap.Boot(products.Blog()...)

	app.Start()
}
