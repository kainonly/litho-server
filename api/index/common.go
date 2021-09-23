package index

import (
	"go.uber.org/fx"
	"lab-api/common"
)

var Provides = fx.Provide(
	NewController,
	NewService,
)

type Controller struct {
	UnimplementedIndexServer

	*ControllerInject
}

type ControllerInject struct {
	common.App

	*Service
}

func NewController(i ControllerInject) *Controller {
	return &Controller{
		ControllerInject: &i,
	}
}

type Service struct {
	*ServiceInject
}

type ServiceInject struct {
	common.App
}

func NewService(i ServiceInject) *Service {
	return &Service{
		ServiceInject: &i,
	}
}
