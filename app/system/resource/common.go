package resource

import (
	"github.com/kainonly/go-bit/crud"
	"github.com/kainonly/go-bit/support"
	"go.uber.org/fx"
	"lab-api/common"
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
		Crud:             crud.New(i.Db, &support.Resource{}),
	}
}

type Service struct {
	*ServiceInject
	Key string
}

type ServiceInject struct {
	common.App
}

func NewService(i ServiceInject) *Service {
	return &Service{
		ServiceInject: &i,
		Key:           i.Set.RedisKey("resource"),
	}
}
