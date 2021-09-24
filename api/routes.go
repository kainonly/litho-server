package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/api"
	"github.com/kainonly/go-bit/mvc"
	"github.com/kainonly/go-bit/passport"
	"go.uber.org/fx"
	"lab-api/api/index"
)

var Options = fx.Options(
	index.Provides,
	fx.Invoke(func(
		route *gin.Engine,
		api *api.API,
		pp *passport.Passport,
		index *index.Controller,
	) {
		mvc.New(route, index)
		mvc.New(route, api, mvc.SetPath("xapi"))
	}),
)
