package orders

import (
	"server/common"

	"github.com/goforj/wire"
)

const (
	Key   = "orders"
	Label = "订单"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	OrdersX *Service
}

type Service struct {
	*common.Inject
}
