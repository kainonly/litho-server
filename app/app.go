package app

import (
	"api/app/pages"
	"api/app/pictures"
	"api/app/roles"
	"api/app/system"
	"api/app/users"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	system.Provides,
	engine.Provides,
	pages.Provides,
	roles.Provides,
	users.Provides,
	pictures.Provides,
	New,
	Subscribe,
)

func New(
	values *common.Values,
	passport *passport.Passport,
	system *system.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	pictures *pictures.Controller,
) *gin.Engine {
	r := globalMiddleware(gin.New(), values)
	r.GET("/", route.Use(system.Index))
	auth := authGuard(passport)

	r.POST("/auth", route.Use(system.AuthLogin))
	r.HEAD("/auth", route.Use(system.AuthVerify))
	r.GET("/auth", auth, route.Use(system.AuthCode))
	r.PUT("/auth", auth, route.Use(system.AuthRefresh))
	r.DELETE("/auth", auth, route.Use(system.AuthLogout))

	r.GET("/uploader", auth, route.Use(system.Uploader))
	r.GET("/navs", auth, route.Use(system.Navs))
	r.GET("/pages/:id", auth, route.Use(system.Dynamic))

	api := r.Group("/api", auth)
	{
		engine.DefaultRouters(api)
		api.PATCH("/:model/sort", auth, route.Use(system.Sort))
		_pages := api.Group("pages")
		{
			_pages.GET("/_indexes/:id", route.Use(pages.GetIndexes))
			_pages.PUT("/_indexes/:id/:index", route.Use(pages.SetIndex))
			_pages.DELETE("/_indexes/:id/:index", route.Use(pages.DeleteIndex))
		}
		_pictures := api.Group("pictures")
		{
			_pictures.GET("/image-info", route.Use(pictures.ImageInfo))
		}
	}
	return r
}
