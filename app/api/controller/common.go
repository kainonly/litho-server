package controller

import (
	"github.com/google/wire"
	"lab-api/app/api/service"
)

type Dependency struct {
	IndexService *service.Index
}

var Provides = wire.NewSet(
	wire.Struct(new(Dependency), "*"),
	NewIndex,
)
