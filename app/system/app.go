package system

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/app/system/controller"
	"lab-api/app/system/service"
)

var App = fx.Options(service.Provides, controller.Provides, fx.Invoke(Routes))

type Dependency struct {
	fx.In

	*controller.Resource
	*controller.Admin
}

func Routes(r *gin.Engine, d Dependency) {
	s := r.Group("/system")
	resource := s.Group("resource")
	{
		resource.POST("originLists", crud.Bind(d.Resource.OriginLists))
		resource.POST("lists", crud.Bind(d.Resource.Lists))
		resource.POST("get", crud.Bind(d.Resource.Get))
		resource.POST("add", crud.Bind(d.Resource.Add))
		resource.POST("edit", crud.Bind(d.Resource.Edit))
		resource.POST("delete", crud.Bind(d.Resource.Delete))
	}
}
