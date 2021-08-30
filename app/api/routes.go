package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kainonly/go-bit/mvc"
	"lab-api/app/api/controller"
	"lab-api/app/api/service"
)

var Provides = wire.NewSet(
	service.Provides,
	controller.Provides,
	wire.Struct(new(Dependency), "*"),
	NewRoutes,
)

type Dependency struct {
	*controller.Index
}

type Routes struct{}

func NewRoutes(r *gin.Engine, d *Dependency) *Routes {
	r.GET("/", mvc.Bind(d.Index.Index))
	return &Routes{}
}
