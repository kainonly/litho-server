package index

import (
	"github.com/kainonly/go-bit/authx"
	"go.uber.org/fx"
	"lab-api/app/system/admin"
	"lab-api/app/system/resource"
	"lab-api/common"
)

var Provides = fx.Provide(
	NewController,
	NewService,
)

type Controller struct {
	*ControllerInject
	Auth *authx.Auth
}

type ControllerInject struct {
	common.App

	Service         *Service
	ResourceService *resource.Service
	AdminService    *admin.Service
}

func NewController(i ControllerInject) *Controller {
	return &Controller{
		ControllerInject: &i,
		Auth:             i.Authx.Make("system"),
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
		Key:           i.Set.RedisKey("code:"),
	}
}
