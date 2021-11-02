package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/support/api"
	"github.com/weplanx/support/helper"
	"github.com/weplanx/support/passport"
	"github.com/weplanx/support/route"
	"laboratory/api/xapi"
	"laboratory/api/xapi/devops"
	"laboratory/api/xapi/schema"
	"laboratory/api/xapi/system"
	"laboratory/common"
	"os"
	"time"
)

var Provides = wire.NewSet(
	xapi.Provides,
	HttpServer,
)

func HttpServer(
	config *common.Set,
	pp *passport.Passport,
	cookie *helper.CookieHelper,
	api *api.API,
	xsystem *system.Controller,
	xdevops *devops.Controller,
	xschema *schema.Controller,
) (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Cors,
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "CONTENT-TYPE"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	x := r.Group("xapi")
	{
		auth := authSystem(pp.Make("system"), cookie)
		route.Auto(x, xsystem, route.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
		if os.Getenv("GIN_MODE") != "release" {
			route.Auto(x, xdevops, route.SetPath("devops"))
		}
		route.Auto(x, xschema, route.SetPath("schema"), route.SetMiddleware(auth))
		route.Auto(x, api, route.SetPath(":collection"), route.SetMiddleware(auth))
	}
	return
}
