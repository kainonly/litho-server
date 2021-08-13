package system

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/app/system/controller"
	"lab-api/app/system/service"
	"lab-api/config"
	"time"
)

var App = fx.Options(service.Provides, controller.Provides, fx.Invoke(Routes))

type Dependency struct {
	fx.In

	*controller.Resource
	*controller.Admin
}

func Routes(r *gin.Engine, d Dependency, config config.Config) {
	s := r.Group("/system")
	cros := config.Cors["system"]
	s.Use(cors.New(cors.Config{
		AllowOrigins:     cros.Origin,
		AllowMethods:     cros.Method,
		AllowHeaders:     cros.AllowHeader,
		ExposeHeaders:    cros.ExposedHeader,
		MaxAge:           time.Duration(cros.MaxAge) * time.Second,
		AllowCredentials: cros.Credentials,
	}))
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
