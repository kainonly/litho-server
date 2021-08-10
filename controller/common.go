package controller

import (
	"github.com/google/wire"
)

type Controllers struct {
	*Index
}

var Provides = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
	NewIndex,
)
