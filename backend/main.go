package main

import (
	"reflexcms/backend/bootstrap"
	"reflexcms/backend/products"
)

func main() {
	app := bootstrap.Boot(products.Full()...)

	app.Start()
}
