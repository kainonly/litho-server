package caps

import (
	"server/common"

	"github.com/goforj/wire"
)

const (
	Key   = "caps"
	Label = "能力标识"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	CapsX *Service
}

type Service struct {
	*common.Inject
}
