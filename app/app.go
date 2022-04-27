package app

import (
	"api/app/departments"
	"api/app/feishu"
	"api/app/pages"
	"api/app/pictures"
	"api/app/roles"
	"api/app/system"
	"api/app/users"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	system.Provides,
	engine.Provides,
	pages.Provides,
	roles.Provides,
	departments.Provides,
	users.Provides,
	feishu.Provides,
	pictures.Provides,
	New,
	Subscribe,
)

func New(
	values *common.Values,
	systemMiddleware *system.Middleware,
	system *system.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	feishu *feishu.Controller,
	pictures *pictures.Controller,
) *gin.Engine {
	r := globalMiddleware(gin.New(), values)
	r.Use(systemMiddleware.RequestLogging())
	auth := systemMiddleware.AuthGuard()

	r.GET("/", route.Use(system.Index))
	r.POST("/auth", route.Use(system.AuthLogin))
	r.HEAD("/auth", route.Use(system.AuthVerify))
	r.GET("/auth", auth, route.Use(system.AuthCode))
	r.PUT("/auth", auth, route.Use(system.AuthRefresh))
	r.DELETE("/auth", auth, route.Use(system.AuthLogout))

	r.GET("/forget-captcha", route.Use(system.ForgetCaptcha))
	r.POST("/forget-verify", route.Use(system.ForgetVerify))
	r.POST("/forget-reset", route.Use(system.ForgetReset))

	r.HEAD("/user/_check", auth, route.Use(system.CheckUser))
	r.GET("/user", auth, route.Use(system.GetUser))
	r.POST("/user", auth, route.Use(system.SetUser))

	r.GET("/vars", auth, route.Use(system.GetVars))
	r.GET("/vars/:key", auth, route.Use(system.GetVar))
	r.PUT("/vars/:key", auth, route.Use(system.SetVar))
	r.GET("/sessions", auth, route.Use(system.GetSessions))
	r.DELETE("/sessions", auth, route.Use(system.DeleteSessions))
	r.DELETE("/sessions/:id", auth, route.Use(system.DeleteSession))

	r.GET("/uploader", auth, route.Use(system.Uploader))
	r.GET("/navs", auth, route.Use(system.Navs))
	r.GET("/pages/:id", auth, route.Use(system.Dynamic))

	_feishu := r.Group("/feishu")
	{
		_feishu.GET("_option", route.Use(feishu.Option))
		_feishu.GET("", route.Use(feishu.OAuth))
		_feishu.POST("", route.Use(feishu.Challenge))
	}

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
