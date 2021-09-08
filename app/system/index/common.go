package index

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit/authx"
	"lab-api/app/system/resource"
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
	Auth *authx.Auth
}

type ControllerInject struct {
	Service         *Service
	ResourceService *resource.Service
}

func NewController(d *common.Dependency, i *ControllerInject, authx *authx.Authx) *Controller {
	return &Controller{
		Dependency:       d,
		ControllerInject: i,
		Auth:             authx.Make("system"),
	}
}

type Service struct {
	*common.Dependency
	Key string
}

func NewService(d *common.Dependency) *Service {
	return &Service{
		Dependency: d,
		Key:        d.App.RedisKey("code:"),
	}
}
