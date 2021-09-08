package resource

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit/crud"
	"github.com/kainonly/go-bit/support"
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
	*crud.Crud
}

type ControllerInject struct {
	Service *Service
}

func NewController(d *common.Dependency, i *ControllerInject) *Controller {
	return &Controller{
		Dependency:       d,
		ControllerInject: i,
		Crud:             crud.New(d.Db, &support.Resource{}),
	}
}

type Service struct {
	*common.Dependency
	Key string
}

func NewService(d *common.Dependency) *Service {
	return &Service{
		Dependency: d,
		Key:        d.Set.RedisKey("resource"),
	}
}
