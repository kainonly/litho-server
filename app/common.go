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
	"os"
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
	roles *roles.Controller,
	departments *departments.Controller,
	users *users.Controller,
	pictures *pictures.Controller,
	videos *videos.Controller,
) *gin.Engine {
	r := globalMiddleware(gin.New(), values)
	r.GET("/", route.Use(index.Index))
	if os.Getenv("GIN_MODE") != "release" {
		r.POST("/install", route.Use(index.Install))
	}
	auth := authGuard(passport)
	r.POST("/auth", route.Use(index.Login))
	r.HEAD("/auth", route.Use(index.Verify))
	r.GET("/auth", auth, route.Use(index.Code))
	r.PUT("/auth", auth, route.Use(index.RefreshToken))
	r.DELETE("/auth", auth, route.Use(index.Logout))
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
		api.POST("/:model", route.Use(engine.Actions))
		api.GET("/:model", route.Use(engine.Find))
		api.GET("/:model/:id", route.Use(engine.FindOneById))
		api.PATCH("/:model", route.Use(engine.Update))
		api.PATCH("/:model/:id", route.Use(engine.UpdateOne))
		api.PUT("/:model/:id", route.Use(engine.ReplaceOne))
		api.DELETE("/:model/:id", route.Use(engine.DeleteOne))
		_pages := api.Group("pages")
		{
			_pages.GET("/:id", route.Use(engine.FindOneById, route.SetModel("pages")))
			_pages.PUT("/:id", route.Use(engine.ReplaceOne, route.SetModel("pages")))
			_pages.DELETE("/:id", route.Use(engine.DeleteOne, route.SetModel("pages")))
			_pages.GET("/has-schema-key", route.Use(pages.HasSchemaKey))
			_pages.PATCH("/sort", route.Use(pages.Sort))
			_pages.GET("/:id/indexes", route.Use(pages.Indexes))
			_pages.PUT("/:id/indexes/:name", route.Use(pages.CreateIndex))
			_pages.DELETE("/:id/indexes/:name", route.Use(pages.DeleteIndex))
		}
		_roles := api.Group("roles")
		{
			_roles.GET("/has-name", route.Use(roles.HasName))
			_roles.GET("/labels", route.Use(roles.Labels))
		}
		_departments := api.Group("departments")
		{
			_departments.PATCH("/sort", route.Use(departments.Sort))
		}
		_users := api.Group("users")
		{
			_users.GET("/has-username", route.Use(users.HasUsername))
			_users.GET("/labels", route.Use(users.Labels))
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
