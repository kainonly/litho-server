package roles

import (
	"server/common"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	RolesX *Service
}

type Service struct {
	*common.Inject
}
