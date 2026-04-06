package products

import (
	"server/common"

	"github.com/goforj/wire"
)

const (
	Key   = "products"
	Label = "产品"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	ProductsX *Service
}

type Service struct {
	*common.Inject
}
