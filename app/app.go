package app

import (
	"api/app/departments"
	"api/app/feishu"
	"api/app/pages"
	"api/app/roles"
	"api/app/system"
	"api/app/tencent"
	"api/app/users"
	"api/app/vars"
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
	tencent.Provides,
	feishu.Provides,
	vars.Provides,
	New,
	Subscribe,
)

func New(
	values *common.Values,
	systemMiddleware *system.Middleware,
	system *system.Controller,
	vars *vars.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	tencent *tencent.Controller,
	feishu *feishu.Controller,
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

	_user := r.Group("/user")
	{
		_user.GET("/captcha", route.Use(system.CaptchaUser))
		_user.POST("/verify", route.Use(system.VerifyUser))
		_user.POST("/reset", route.Use(system.ResetUser))
		_user.HEAD("", auth, route.Use(system.CheckUser))
		_user.GET("", auth, route.Use(system.GetUser))
		_user.POST("", auth, route.Use(system.SetUser))
	}

	_vars := r.Group("/vars")
	{
		_vars.GET("", auth, route.Use(vars.Gets))
		_vars.GET("/_option", auth, route.Use(vars.Option))
		_vars.GET("/:key", auth, route.Use(vars.Get))
		_vars.PUT("/:key", auth, route.Use(vars.Set))
	}

	r.GET("/sessions", auth, route.Use(system.GetSessions))
	r.DELETE("/sessions", auth, route.Use(system.DeleteSessions))
	r.DELETE("/sessions/:id", auth, route.Use(system.DeleteSession))

	r.GET("/upload", auth, route.Use(system.Upload))

	r.GET("/navs", auth, route.Use(system.Navs))
	r.GET("/pages/:id", auth, route.Use(system.Dynamic))

	_tencent := r.Group("/tencent")
	{
		_tencent.GET("cos/presigned", auth, route.Use(tencent.CosPresigned))
		_tencent.GET("cos/image-info", auth, route.Use(tencent.ImageInfo))
	}

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
	}
	return r
}
