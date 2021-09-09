package admin

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit/crud"
	"lab-api/common"
	"lab-api/model"
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
		Crud:             crud.New(d.Db, &model.Admin{}),
	}
}

type Service struct {
	*common.Dependency
	Key string
}

func NewService(d *common.Dependency) *Service {
	return &Service{
		Dependency: d,
		Key:        d.Set.RedisKey("admin"),
	}
}
