package index

import (
	"github.com/kainonly/go-bit/authx"
	"go.uber.org/fx"
	"lab-api/app/system/resource"
	"lab-api/common"
)

var Provides = fx.Provide(
	NewController,
	NewService,
)

type Controller struct {
	*common.Dependency
	*ControllerInject
	Auth *authx.Auth
}

type ControllerInject struct {
	fx.In

	Service         *Service
	ResourceService *resource.Service
}

func NewController(d common.Dependency, i ControllerInject, authx *authx.Authx) *Controller {
	return &Controller{
		Dependency:       &d,
		ControllerInject: &i,
		Auth:             authx.Make("system"),
	}
}

type Service struct {
	*common.Dependency
	Key string
}

func NewService(d common.Dependency) *Service {
	return &Service{
		Dependency: &d,
		Key:        d.App.RedisKey("code:"),
	}
}
