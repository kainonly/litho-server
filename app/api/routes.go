package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kainonly/go-bit/mvc"
	"lab-api/app/api/dev"
	"lab-api/app/api/index"
)

var Provides = wire.NewSet(
	wire.Struct(new(Inject), "*"),
	index.Provides,
	dev.Provides,
	NewRoutes,
)

type Inject struct {
	Index *index.Controller
	Dev   *dev.Controller
}

type Routes struct{}

func NewRoutes(r *gin.Engine, i *Inject) *Routes {
	r.GET("/", mvc.Bind(i.Index.Index))
	sub := r.Group("/dev")
	{
		sub.POST("setup", mvc.Bind(i.Dev.Setup))
		sub.POST("sync", mvc.Bind(i.Dev.Sync))
		sub.POST("migrate", mvc.Bind(i.Dev.Migrate))
	}
	return &Routes{}
}
