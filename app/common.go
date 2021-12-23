package app

import (
	"api/app/index"
	"api/app/pages"
	"api/app/roles"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/api"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/route"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	pages.Provides,
	roles.Provides,
	fx.Invoke(func(
		r *gin.Engine,
		api *api.Controller,
		passport *passport.Passport,
		index *index.Controller,
	) {
		r.GET("/", route.Use(Home))
	}),
)
