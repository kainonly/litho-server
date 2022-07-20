package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/server/api/app"
	"github.com/weplanx/server/api/dsl"
	"github.com/weplanx/server/api/values"
)

var Provides = wire.NewSet(
	wire.Struct(new(API), "*"),
	app.Provides,
	values.Provides,
	dsl.Provides,
	Routes,
)

func Routes(
	api *API,
	app *app.Controller,
	dsl *dsl.Controller,
) *gin.Engine {
	r := api.Engine()

	app.In(r.Group("/"))

	dsl.In(r.Group("/dsl/:model"))

	return r
}
