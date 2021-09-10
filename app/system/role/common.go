package role

import (
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/common"
	"lab-api/model"
)

var Provides = fx.Provide(
	NewController,
	NewService,
)

type Controller struct {
	*ControllerInject
	*crud.Crud
}

type ControllerInject struct {
	common.App

	Service *Service
}

func NewController(i ControllerInject) *Controller {
	return &Controller{
		ControllerInject: &i,
		Crud:             crud.New(i.Db, &model.Role{}),
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
