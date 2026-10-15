package tencent

import (
	"litho-api/common"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Service), "*"),
)

type Service struct {
	*common.Inject
}
