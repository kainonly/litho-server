package routes

import (
	"server/common"

	"github.com/goforj/wire"
)

const (
	Resource = "/routes"
	Label    = "路由"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	RoutesX *Service
}
type Service struct {
	*common.Inject
}
