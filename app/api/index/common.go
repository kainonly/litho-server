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
	*common.Dependency
	*Service
}

func NewController(d common.Dependency, service *Service) *Controller {
	return &Controller{
		Dependency: &d,
		Service:    service,
	}
}

type Service struct {
	*common.Dependency
}

func NewService(d common.Dependency) *Service {
	return &Service{
		Dependency: &d,
	}
}
