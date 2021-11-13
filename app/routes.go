package app

import (
	"api/app/index"
	"api/app/x"
	"api/app/x/devops"
	"api/app/x/schema"
	"api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
	"os"
	"time"
)

var Provides = wire.NewSet(
	index.Provides,
	x.Provides,
	HttpServer,
)

func HttpServer(
	config *common.Set,
	pp *passport.Passport,
	cookie *helper.CookieHelper,
	index *index.Controller,
	api *api.API,
	xindex *x.Controller,
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
	r.GET("/", route.Returns(index.Index))
	xrg := r.Group("x")
	{
		auth := authSystem(pp.Make("system"), cookie)
		route.Auto(xrg, xindex, route.SetMiddleware(auth, "Code", "RefreshToken", "Logout", "Pages"))
		if os.Getenv("GIN_MODE") != "release" {
			route.Auto(xrg, xdevops, route.SetPath("devops"))
		}
		route.Auto(xrg, xschema, route.SetPath("schema"), route.SetMiddleware(auth))
		route.Auto(xrg, api, route.SetPath(":collection"), route.SetMiddleware(auth))
	}
	return
}
