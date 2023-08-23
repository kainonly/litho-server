package schedules

import (
	"github.com/google/wire"
	"github.com/weplanx/server/common"
	"sync"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	NewService,
)

func NewService(i *common.Inject) *Service {
	return &Service{
		Inject: i,
		M:      sync.Map{},
	}
}

type M map[string]interface{}
