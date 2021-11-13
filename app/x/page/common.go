package page

import (
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(InjectController), "*"),
	wire.Struct(new(InjectService), "*"),
	NewController,
	NewService,
)
