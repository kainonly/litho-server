package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/support/api"
	"github.com/weplanx/support/passport"
	"github.com/weplanx/support/route"
	"go.uber.org/fx"
	"laboratory/api/xapi"
	"laboratory/api/xapi/devops"
	"laboratory/api/xapi/system"
)

var Options = fx.Options(
	xapi.Provides,
	fx.Invoke(func(
		r *gin.Engine,
		pp *passport.Passport,
		api *api.API,
		xsystem *system.Controller,
		xdevops *devops.Controller,
	) {
		xapi := r.Group("xapi")
		{
			route.Auto(xapi, xsystem)
			route.Auto(xapi, xdevops, route.SetPath("devops"))
			route.Auto(xapi, api, route.SetPath(":model"))
		}
	}),
)
