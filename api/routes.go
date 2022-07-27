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
	"github.com/weplanx/server/api/sessions"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/api/values"
)

var Provides = wire.NewSet(
	wire.Struct(new(API), "*"),
	index.Provides,
	values.Provides,
	sessions.Provides,
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
	sessions *sessions.Controller,
	dsl *dsl.Controller,
) (h *server.Hertz, err error) {
	if h, err = api.Engine(); err != nil {
		return
	}
	var auth *jwt.HertzJWTMiddleware
	if auth, err = api.Auth(); err != nil {
		return
	}

	h.POST("login", auth.LoginHandler)

	app := h.Group("", auth.MiddlewareFunc())
	{
		app.GET("", index.Index)
		app.GET("code", index.GetRefreshCode)
		app.POST("refresh_token", index.VerifyRefreshCode, auth.RefreshHandler)

		app.GET("user", index.GetUser)
		app.PATCH("user", index.SetUser)
		app.DELETE("user", auth.LogoutHandler)

		values.In(app.Group("values"))
		sessions.In(app.Group("sessions"))
		dsl.In(app.Group("dsl/:model"))
	}

	return
}
