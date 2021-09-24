package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/api/index"
	"lab-api/common"
)

var Options = fx.Options(index.Provides, fx.Invoke(Routes))

type Inject struct {
	common.App

	Index *index.Controller
}

func Routes(r *gin.Engine, i Inject) {
	mvc.New(r, i.Index)
	mvc.New(r, i.API, mvc.SetPath("dsapi"))
}
