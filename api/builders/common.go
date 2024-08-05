package builders

import (
	"github.com/google/wire"
	"github.com/weplanx/server/common"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	BuildersX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]interface{}
