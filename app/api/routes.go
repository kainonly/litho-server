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
}

func Routes(r *gin.Engine, d Dependency) {
	r.GET("/", mvc.Bind(d.Index.Index))
	r.POST("/setup", mvc.Bind(d.Index.Setup))
	r.PUT("/sync", mvc.Bind(d.Index.Sync))
}
