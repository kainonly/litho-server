package app

import (
	"api/app/departments"
	"api/app/feishu"
	"api/app/pages"
	"api/app/roles"
	"api/app/sessions"
	"api/app/tencent"
	"api/app/user"
	"api/app/users"
	"api/app/vars"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	wire.Struct(new(Middleware), "*"),
	user.Provides,
	engine.Provides,
	tencent.Provides,
	feishu.Provides,
	sessions.Provides,
	vars.Provides,
	pages.Provides,
	roles.Provides,
	departments.Provides,
	users.Provides,
	New,
	Subscribe,
)

func New(
	middleware *Middleware,
	user *user.Controller,
	vars *vars.Controller,
	sessions sessions.Controller,
	tencent *tencent.Controller,
	feishu *feishu.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
) *gin.Engine {
	r := middleware.Global()
	auth := middleware.AuthGuard()

	r.POST("/auth", route.Use(user.AuthLogin))
	r.HEAD("/auth", route.Use(user.AuthVerify))
	r.GET("/auth", auth, route.Use(user.AuthCode))
	r.PUT("/auth", auth, route.Use(user.AuthRefresh))
	r.DELETE("/auth", auth, route.Use(user.AuthLogout))

	_user := r.Group("/user")
	{
		_user.GET("/captcha", route.Use(user.GetCaptcha))
		_user.POST("/captcha", route.Use(user.VerifyCaptcha))
		_user.POST("/reset", route.Use(user.Reset))
		_user.HEAD("_exists", auth, route.Use(user.Exists))
		_user.GET("", auth, route.Use(user.Get))
		_user.POST("", auth, route.Use(user.Set))
	}

	r.GET("options", auth, route.Use(vars.Options))
	_vars := r.Group("/vars", auth)
	{
		_vars.GET("", route.Use(vars.Gets))
		_vars.GET("/:key", route.Use(vars.Get))
		_vars.PUT("/:key", route.Use(vars.Set))
	}

	_sessions := r.Group("/sessions", auth)
	{
		_sessions.GET("", route.Use(sessions.Gets))
		_sessions.DELETE("", route.Use(sessions.BulkDelete))
		_sessions.DELETE("/:id", route.Use(sessions.Delete))
	}

	_tencent := r.Group("/tencent", auth)
	{
		_tencent.GET("cos-presigned", route.Use(tencent.CosPresigned))
		_tencent.GET("cos-image-info", route.Use(tencent.ImageInfo))
	}

	_feishu := r.Group("/feishu")
	{
		_feishu.POST("", route.Use(feishu.Challenge))
		_feishu.GET("", route.Use(feishu.OAuth))
		_feishu.GET("option", auth, route.Use(feishu.Option))
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
