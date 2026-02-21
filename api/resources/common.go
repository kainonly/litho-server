package resources

import (
	"server/common"

	"github.com/goforj/wire"
)

const (
	Key = "resources"
	Label    = "资源"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	ResourcesX *Service
}

type Service struct {
	*common.Inject
}
