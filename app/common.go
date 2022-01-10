package app

import (
	"api/app/index"
	"api/app/pages"
	"api/app/roles"
	"api/app/users"
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
	pages.Provides,
	roles.Provides,
	users.Provides,
	New,
)

func New(
	values *common.Values,
	passport *passport.Passport,
	index *index.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	roles *roles.Controller,
	users *users.Controller,
) *gin.Engine {
	r := globalMiddleware(gin.New(), values)
	r.GET("/", route.Use(index.Index))
	auth := authGuard(passport)
	r.POST("/auth", route.Use(index.Login))
	r.HEAD("/auth", route.Use(index.Verify))
	r.GET("/auth", auth, route.Use(index.Code))
	r.PUT("/auth", auth, route.Use(index.RefreshToken))
	r.DELETE("/auth", auth, route.Use(index.Logout))
	api := r.Group("/api", auth)
	{
		api.GET("", route.Use(index.Api))
		api.POST("/:model", route.Use(engine.Create))
		api.GET("/:model", route.Use(engine.Find))
		api.GET("/:model/:id", route.Use(engine.FindOneById))
		api.PATCH("/:model", route.Use(engine.Update))
		api.PATCH("/:model/:id", route.Use(engine.UpdateOneById))
		api.PUT("/:model/:id", route.Use(engine.ReplaceOneById))
		api.DELETE("/:model/:id", route.Use(engine.DeleteOneById))
		_pages := api.Group("pages")
		{
			_pages.GET("/:id", route.Use(engine.FindOneById, route.SetModel("pages")))
			_pages.PUT("/:id", route.Use(engine.ReplaceOneById, route.SetModel("pages")))
			_pages.DELETE("/:id", route.Use(engine.DeleteOneById, route.SetModel("pages")))
			_pages.GET("/has-schema-key", route.Use(pages.HasSchemaKey))
			_pages.PATCH("/sort", route.Use(pages.Sort))
			_pages.GET("/:id/indexes", route.Use(pages.FindIndexes))
			_pages.PUT("/:id/indexes/:name", route.Use(pages.CreateIndex))
			_pages.DELETE("/:id/indexes/:name", route.Use(pages.DeleteIndex))
		}
	}
	return r
}
