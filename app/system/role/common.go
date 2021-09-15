package role

import (
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/common"
)

var Provides = fx.Provide(
	NewController,
	NewService,
)

type Controller struct {
	*ControllerInject
	*crud.API
}

type ControllerInject struct {
	common.App

	Service *Service
}

func NewController(i ControllerInject) *Controller {
	return &Controller{
		ControllerInject: &i,
		API:              i.Crud.API("role"),
	}
}

type Service struct {
	*ServiceInject
	Routers     string
	Permissions string
}

type ServiceInject struct {
	common.App
}

func NewService(i ServiceInject) *Service {
	return &Service{
		ServiceInject: &i,
		Routers:       i.Set.RedisKey("routers"),
		Permissions:   i.Set.RedisKey("permissions"),
	}
}
