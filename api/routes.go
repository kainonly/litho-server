package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/support/api"
	"github.com/weplanx/support/helper"
	"github.com/weplanx/support/passport"
	"github.com/weplanx/support/route"
	"go.uber.org/fx"
	"laboratory/api/xapi"
	"laboratory/api/xapi/devops"
	"laboratory/api/xapi/schema"
	"laboratory/api/xapi/system"
	"os"
)

var Options = fx.Options(
	xapi.Provides,
	fx.Invoke(func(
		r *gin.Engine,
		pp *passport.Passport,
		cookie *helper.CookieHelper,
		api *api.API,
		xsystem *system.Controller,
		xdevops *devops.Controller,
		xschema *schema.Controller,
	) {
		xapi := r.Group("xapi")
		{
			auth := authSystem(pp.Make("system"), cookie)
			route.Auto(xapi, xsystem, route.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
			if os.Getenv("GIN_MODE") != "release" {
				route.Auto(xapi, xdevops, route.SetPath("devops"))
			}
			route.Auto(xapi, xschema, route.SetPath("schema"), route.SetMiddleware(auth))
			route.Auto(xapi, api, route.SetPath(":collection"), route.SetMiddleware(auth))
		}
	}),
)
