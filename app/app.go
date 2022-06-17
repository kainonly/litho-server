package app

import (
	"api/app/departments"
	"api/app/feishu"
	"api/app/pages"
	"api/app/roles"
	"api/app/schedules"
	"api/app/system"
	"api/app/tencent"
	"api/app/users"
	"api/app/values"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	wire.Struct(new(Middleware), "*"),
	values.Provides,
	schedules.Provides,
	system.Provides,
	engine.Provides,
	tencent.Provides,
	feishu.Provides,
	pages.Provides,
	roles.Provides,
	departments.Provides,
	users.Provides,
	New,
	SetJobs,
)

func New(
	middleware *Middleware,
	system *system.Controller,
	values *values.Controller,
	schedules *schedules.Controller,
	tencent *tencent.Controller,
	feishu *feishu.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	_ *common.Jobs,
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
	r.DELETE("/values/:key", auth, route.Use(values.Delete))

	r.GET("/navs", auth, route.Use(pages.Navs))

	_pages := r.Group("pages", auth)
	{
		// TODO: 未使用 DSL 的原因是要适配为缓存，暂不处理
		_pages.GET("/:id", route.Use(pages.Dynamic))

		// 索引管理
		_pages.GET("/indexes/:id", route.Use(pages.GetIndexes))
		_pages.PUT("/indexes/:id/:index", route.Use(pages.SetIndex))
		_pages.DELETE("/indexes/:id/:index", route.Use(pages.DeleteIndex))
	}

	_schedules := r.Group("schedules")
	{
		_schedules.GET("/", route.Use(schedules.List))
		_schedules.GET("/:key", route.Use(schedules.Get))
		_schedules.POST("/sync", route.Use(schedules.SetSync))
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
	}

	// 设置 Query DSL 路由
	engine.SetRouters(r.Group("/dsl", auth))
	return r
}
