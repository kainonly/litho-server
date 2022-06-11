package app

import (
	"api/app/departments"
	"api/app/feishu"
	"api/app/pages"
	"api/app/roles"
	"api/app/system"
	"api/app/tencent"
	"api/app/users"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/route"
	"github.com/weplanx/go/values"
)

var Provides = wire.NewSet(
	wire.Struct(new(Middleware), "*"),
	system.Provides,
	engine.Provides,
	tencent.Provides,
	feishu.Provides,
	pages.Provides,
	roles.Provides,
	departments.Provides,
	users.Provides,
	values.Provides,
	New,
	Subscribe,
)

func New(
	middleware *Middleware,
	system *system.Controller,
	tencent *tencent.Controller,
	feishu *feishu.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	values *values.Controller,
) *gin.Engine {
	r := middleware.Global()
	auth := middleware.AuthGuard()

	r.GET("/", route.Use(system.Index))

	r.POST("/auth", route.Use(system.AuthLogin))
	r.HEAD("/auth", route.Use(system.AuthVerify))
	r.GET("/auth", auth, route.Use(system.AuthCode))
	r.PUT("/auth", auth, route.Use(system.AuthRefresh))
	r.DELETE("/auth", auth, route.Use(system.AuthLogout))
	r.GET("/captcha", route.Use(system.GetCaptcha))
	r.POST("/captcha", route.Use(system.VerifyCaptcha))
	r.HEAD("/user/_exists", auth, route.Use(system.ExistsUser))
	r.GET("/user", auth, route.Use(system.GetUser))
	r.POST("/user", auth, route.Use(system.SetUser))
	r.POST("/user/reset", route.Use(system.ResetUser))
	r.GET("/sessions", auth, route.Use(system.GetSessions))
	r.DELETE("/sessions", auth, route.Use(system.DeleteSessions))
	r.DELETE("/sessions/:id", auth, route.Use(system.DeleteSession))

	r.GET("/options", route.Use(system.Options))
	r.GET("/values", auth, route.Use(values.Get))
	r.PATCH("/values", auth, route.Use(values.Set))

	_tencent := r.Group("/tencent", auth)
	{
		_tencent.GET("cos-presigned", route.Use(tencent.CosPresigned))
		_tencent.GET("cos-image-info", route.Use(tencent.ImageInfo))
	}

	_feishu := r.Group("/feishu")
	{
		_feishu.POST("", route.Use(feishu.Challenge))
		_feishu.GET("", route.Use(feishu.OAuth))
	}

	r.GET("/navs", auth, route.Use(pages.Navs))
	r.GET("/pages/:id", auth, route.Use(pages.Dynamic))

	api := r.Group("/api", auth)
	{
		engine.DefaultRouters(api)
		_pages := api.Group("pages")
		{
			_pages.GET("/_indexes/:id", route.Use(pages.GetIndexes))
			_pages.PUT("/_indexes/:id/:index", route.Use(pages.SetIndex))
			_pages.DELETE("/_indexes/:id/:index", route.Use(pages.DeleteIndex))
		}
	}
	return r
}
