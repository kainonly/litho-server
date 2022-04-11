package app

import (
	"api/app/center"
	"api/app/departments"
	"api/app/index"
	"api/app/pages"
	"api/app/pictures"
	"api/app/roles"
	"api/app/users"
	"api/app/videos"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	index.Provides,
	engine.Provides,
	center.Provides,
	pages.Provides,
	roles.Provides,
	departments.Provides,
	users.Provides,
	pictures.Provides,
	videos.Provides,
	New,
	Subscribe,
)

func New(
	values *common.Values,
	passport *passport.Passport,
	index *index.Controller,
	engine *engine.Controller,
	center *center.Controller,
	pages *pages.Controller,
	departments *departments.Controller,
	pictures *pictures.Controller,
	videos *videos.Controller,
) *gin.Engine {
	r := globalMiddleware(gin.New(), values)
	r.GET("/", route.Use(index.Index))
	auth := authGuard(passport)

	r.POST("/auth", route.Use(index.AuthLogin))
	r.HEAD("/auth", route.Use(index.AuthVerify))
	r.GET("/auth", auth, route.Use(index.AuthCode))
	r.PUT("/auth", auth, route.Use(index.AuthRefresh))
	r.DELETE("/auth", auth, route.Use(index.AuthLogout))

	r.GET("/uploader", auth, route.Use(index.Uploader))
	r.GET("/navs", auth, route.Use(index.Navs))
	r.GET("/pages/:id", auth, route.Use(index.Dynamic))
	_center := r.Group("/center", auth)
	{
		_center.GET("/user-info", route.Use(center.GetUserInfo))
		_center.PATCH("/user-info", route.Use(center.SetUserInfo))
	}
	api := r.Group("/api", auth)
	{
		engine.DefaultRouters(api)
		_pages := api.Group("pages")
		{
			_pages.GET("/:id", route.Use(engine.Get, route.SetModel("pages")))
			_pages.PUT("/:id", route.Use(engine.Put, route.SetModel("pages")))
			_pages.DELETE("/:id", route.Use(engine.Delete, route.SetModel("pages")))
			_pages.GET("/has-schema-key", route.Use(pages.HasSchemaKey))
			_pages.PATCH("/sort", route.Use(pages.Sort))
			_pages.GET("/:id/indexes", route.Use(pages.Indexes))
			_pages.PUT("/:id/indexes/:name", route.Use(pages.CreateIndex))
			_pages.DELETE("/:id/indexes/:name", route.Use(pages.DeleteIndex))
		}
		_departments := api.Group("departments")
		{
			_departments.PATCH("/sort", route.Use(departments.Sort))
		}
		_pictures := api.Group("pictures")
		{
			_pictures.GET("/image-info", route.Use(pictures.ImageInfo))
			_pictures.GET("/labels", route.Use(pictures.FindLabels))
			_pictures.POST("/bulk-delete", route.Use(pictures.BulkDelete))
		}
		_videos := api.Group("videos")
		{
			_videos.GET("/labels", route.Use(videos.FindLabels))
			_videos.POST("/bulk-delete", route.Use(videos.BulkDelete))
		}
	}
	return r
}
