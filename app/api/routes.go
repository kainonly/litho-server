package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/app/api/controller"
	"lab-api/app/api/service"
)

var App = fx.Options(service.Provides, controller.Provides, fx.Invoke(Routes))

type Dependency struct {
	fx.In

	*controller.Index
	*controller.Developer
}

func Routes(r *gin.Engine, d Dependency) {
	r.GET("/", mvc.Bind(d.Index.Index))
	dev := r.Group("/dev")
	{
		dev.POST("setup", mvc.Bind(d.Developer.Setup))
		dev.POST("sync", mvc.Bind(d.Developer.Sync))
		dev.POST("migrate", mvc.Bind(d.Developer.Migrate))
	}
}
