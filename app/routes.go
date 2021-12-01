package app

import (
	"api/app/x"
	"api/app/x/devops"
	"api/app/x/schema"
	"github.com/gin-gonic/gin"
	wpx "github.com/weplanx/go"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"go.uber.org/fx"
	"os"
)

var Options = fx.Options(
	x.Options,
	fx.Invoke(func(
		r *gin.Engine,
		pp *passport.Passport,
		cookie *helper.CookieHelper,
		api *api.API,
		x *x.Controller,
		xdevops *devops.Controller,
		xschema *schema.Controller,
	) {
		xapi := r.Group("x")
		{
			auth := authSystem(pp.Make("system"), cookie)
			wpx.Auto(xapi, x, wpx.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
			if os.Getenv("GIN_MODE") != "release" {
				wpx.Auto(xapi, xdevops, wpx.SetPath("devops"))
			}
			wpx.Auto(xapi, xschema, wpx.SetPath("schema"), wpx.SetMiddleware(auth))
			wpx.Auto(xapi, api, wpx.SetPath(":collection"), wpx.SetMiddleware(auth))
		}
	}),
)
