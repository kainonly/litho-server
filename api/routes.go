package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/api"
	"github.com/kainonly/go-bit/mvc"
	"github.com/kainonly/go-bit/passport"
	"go.uber.org/fx"
	"lab-api/api/xapi"
	"lab-api/api/xapi/devops"
	"lab-api/api/xapi/system"
)

var Options = fx.Options(
	xapi.Provides,
	fx.Invoke(func(
		route *gin.Engine,
		pp *passport.Passport,
		api *api.API,
		xsystem *system.Controller,
		xdevops *devops.Controller,
	) {
		xapi := route.Group("xapi")
		{
			mvc.New(xapi, xsystem)
			mvc.New(xapi, xdevops, mvc.SetPath("devops"))
			mvc.New(xapi, api, mvc.SetPath(":model"))
		}

	}),
)
