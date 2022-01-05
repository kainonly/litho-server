package app

import (
	"api/app/index"
	"api/app/pages"
	"api/app/roles"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/engine"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	index.Provides,
	engine.Provides,
	pages.Provides,
	roles.Provides,
	New,
)

func New(
	values *common.Values,
	index *index.Controller,
	engine *engine.Controller,
	pages *pages.Controller,
	roles *roles.Controller,
) *gin.Engine {
	r := middleware(gin.Default(), values)
	r.GET("/", route.Use(index.Index))
	r.POST("/:model", route.Use(engine.Create))
	r.GET("/:model", route.Use(engine.Find))
	r.GET("/:model/:id", route.Use(engine.FindOneById))
	r.PATCH("/:model", route.Use(engine.Update))
	r.PATCH("/:model/:id", route.Use(engine.UpdateOneById))
	r.PUT("/:model/:id", route.Use(engine.ReplaceOneById))
	r.DELETE("/:model/:id", route.Use(engine.DeleteOneById))
	_pages := r.Group("pages")
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
	return r
}
