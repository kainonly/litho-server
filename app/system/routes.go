package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"lab-api/app/system/controller"
	"lab-api/app/system/service"
	"lab-api/common"
)

var Provides = wire.NewSet(
	service.Provides,
	controller.Provides,
	wire.Struct(new(Dependency), "*"),
	NewRoutes,
)

type Dependency struct {
	Config common.Config

	*controller.Index
}

type Routes struct{}

func NewRoutes(r *gin.Engine, d *Dependency) *Routes {
	s := r.Group("/system")
	s.GET("/")
	return &Routes{}
}
