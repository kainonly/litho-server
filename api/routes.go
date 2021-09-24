package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/dsapi"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/api/index"
)

var Options = fx.Options(
	index.Provides,
	fx.Invoke(func(
		route *gin.Engine,
		api *dsapi.API,
		authx *authx.Authx,
		index *index.Controller,
	) {
		mvc.New(route, index)
		mvc.New(route, api, mvc.SetPath("dsapi"))
	}),
)
