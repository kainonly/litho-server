package app

import (
	"api/app/index"
	"api/app/pages"
	"api/app/roles"
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/route"
)

var Provides = wire.NewSet(
	index.Provides,
	api.Provides,
	pages.Provides,
	roles.Provides,
	New,
)

func New(
	values *common.Values,
	index *index.Controller,
	api *api.Controller,
	pages *pages.Controller,
	roles *roles.Controller,
) *gin.Engine {
	r := middleware(gin.Default(), values)
	r.GET("/", route.Use(index.Index))
	r.POST("/:model", route.Use(api.Create))
	r.GET("/:model", route.Use(api.Find))
	r.GET("/:model/:id", route.Use(api.FindOneById))
	r.PATCH("/:model", route.Use(api.Update))
	r.PATCH("/:model/:id", route.Use(api.UpdateOneById))
	r.PUT("/:model/:id", route.Use(api.ReplaceOneById))
	r.DELETE("/:model/:id", route.Use(api.DeleteOneById))
	_pages := r.Group("pages")
	{
		_pages.GET("/:id", route.Use(api.FindOneById, route.SetModel("pages")))
		_pages.PUT("/:id", route.Use(api.ReplaceOneById, route.SetModel("pages")))
		_pages.DELETE("/:id", route.Use(api.DeleteOneById, route.SetModel("pages")))
		_pages.GET("/has-schema-key", route.Use(pages.HasSchemaKey))
		_pages.PATCH("/sort", route.Use(pages.Sort))
		_pages.GET("/:id/indexes", route.Use(pages.FindIndexes))
		_pages.PUT("/:id/indexes/:name", route.Use(pages.CreateIndex))
		_pages.DELETE("/:id/indexes/:name", route.Use(pages.DeleteIndex))
	}
	return r
}
