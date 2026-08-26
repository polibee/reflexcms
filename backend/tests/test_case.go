package tests

import (
	"github.com/goravel/framework/testing"

	"reflexcms/backend/bootstrap"
	"reflexcms/backend/products"
)

func init() {
	bootstrap.Boot(products.Full()...)
}

type TestCase struct {
	testing.TestCase
}
