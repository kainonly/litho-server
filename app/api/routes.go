package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/app/api/index"
	"lab-api/common"
)

var Options = fx.Options(index.Provides, fx.Invoke(Routes))

type Inject struct {
	common.App

	Index *index.Controller
}

func Routes(r *gin.Engine, i Inject) {
	r.GET("/", mvc.Returns(i.Index.Index))
}
