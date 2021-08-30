package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kainonly/go-bit/mvc"
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
	*controller.Admin
}

type Routes struct{}

func NewRoutes(r *gin.Engine, d *Dependency) *Routes {
	s := r.Group("/system")
	admin := s.Group("admin")
	{
		admin.POST("originLists", mvc.Bind(d.Admin.OriginLists))
		admin.POST("lists", mvc.Bind(d.Admin.Lists))
		admin.POST("get", mvc.Bind(d.Admin.Get))
		admin.POST("add", mvc.Bind(d.Admin.Add))
		admin.POST("edit", mvc.Bind(d.Admin.Edit))
		admin.POST("delete", mvc.Bind(d.Admin.Delete))
	}
	return &Routes{}
}
