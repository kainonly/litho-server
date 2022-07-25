package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/hertz-contrib/jwt"
	"github.com/weplanx/server/api/departments"
	"github.com/weplanx/server/api/dsl"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/pages"
	"github.com/weplanx/server/api/roles"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/api/values"
)

var Provides = wire.NewSet(
	wire.Struct(new(API), "*"),
	index.Provides,
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
	index *index.Controller,
	values *values.Controller,
	dsl *dsl.Controller,
) (h *server.Hertz, err error) {
	if h, err = api.Engine(); err != nil {
		return
	}
	var auth *jwt.HertzJWTMiddleware
	if auth, err = api.Auth(); err != nil {
		return
	}

	h.GET("/", index.Index)
	h.POST("login", auth.LoginHandler)

	app := h.Group("", auth.MiddlewareFunc())
	{
		app.DELETE("auth", auth.LogoutHandler)
		app.GET("refresh_code")
		app.GET("refresh_token", auth.RefreshHandler)
		app.GET("user", index.GetUser)
		app.PATCH("user", index.SetUser)

		values.In(app.Group("values"))
		dsl.In(app.Group("dsl/:model"))
	}

	return
}
