package controller

import (
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
	NewIndex,
)

type Controllers struct {
	*Index
}
