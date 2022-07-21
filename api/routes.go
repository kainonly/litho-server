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
	values *values.Controller,
	dsl *dsl.Controller,
) (r *gin.Engine, err error) {
	if r, err = api.Engine(); err != nil {
		return
	}

	app.In(r.Group(""))
	values.In(r.Group("values"))
	dsl.In(r.Group("dsl/:model"))

	return
}
