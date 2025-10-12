package schedulers

import (
	"server/common"

	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	V *common.Values

	SchedulesX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]any
