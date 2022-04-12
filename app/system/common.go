package system

import (
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Middleware), "*"),
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)
