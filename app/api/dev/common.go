package dev

import (
	"github.com/google/wire"
	"lab-api/common"
)

var Provides = wire.NewSet(
	wire.Struct(new(ControllerInject), "*"),
	NewController,
	NewService,
)

type Controller struct {
	*common.Dependency
	*ControllerInject
}

type ControllerInject struct {
	*Service
}

func NewController(d *common.Dependency, i *ControllerInject) *Controller {
	return &Controller{
		Dependency:       d,
		ControllerInject: i,
	}
}

type Service struct {
	*common.Dependency
}

func NewService(d *common.Dependency) *Service {
	return &Service{
		Dependency: d,
	}
}
