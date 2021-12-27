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
	r := middleware(gin.New(), values)
	r.GET("/", route.Use(index.Index))

	api.Auto(r)
	_pages := r.Group("pages")
	{
		_pages.GET("has-schema-key", route.Use(pages.HasSchemaKey))
		_pages.PATCH("sort", route.Use(pages.Sort))
		_pages.GET(":id/indexes", route.Use(pages.FindIndexes))
		_pages.PUT(":id/indexes/:name", route.Use(pages.CreateIndex))
		_pages.DELETE(":id/indexes/:name", route.Use(pages.DeleteIndex))
	}

	return r
}
