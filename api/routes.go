package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/weplanx/server/api/app"
	"github.com/weplanx/server/api/departments"
	"github.com/weplanx/server/api/dsl"
	"github.com/weplanx/server/api/pages"
	"github.com/weplanx/server/api/roles"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/api/values"
)

var Provides = wire.NewSet(
	wire.Struct(new(API), "*"),
	app.Provides,
	values.Provides,
	dsl.Provides,
	pages.Provides,
	users.Provides,
	roles.Provides,
	departments.Provides,
	Routes,
)

func Routes(
	api *API,
	app *app.Controller,
	values *values.Controller,
	dsl *dsl.Controller,
) (h *server.Hertz, err error) {
	if h, err = api.Engine(); err != nil {
		return
	}

	app.In(h.Group(""))
	values.In(h.Group("values"))
	dsl.In(h.Group("dsl/:model"))

	return
}
