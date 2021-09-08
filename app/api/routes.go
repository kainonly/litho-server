package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/app/api/dev"
	"lab-api/app/api/index"
)

var App = fx.Options(index.Provides, dev.Provides, fx.Invoke(Routes))

type Inject struct {
	fx.In

	Index *index.Controller
	Dev   *dev.Controller
}

func Routes(r *gin.Engine, i Inject) {
	r.GET("/", mvc.Bind(i.Index.Index))
	sub := r.Group("/dev")
	{
		sub.POST("setup", mvc.Bind(i.Dev.Setup))
		sub.POST("sync", mvc.Bind(i.Dev.Sync))
		sub.POST("migrate", mvc.Bind(i.Dev.Migrate))
	}
}
