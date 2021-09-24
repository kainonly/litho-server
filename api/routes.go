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
	r.GET("/", mvc.Returns(i.Index.Index))
	r.POST("dsapi/find-one", mvc.Returns(i.API.FindOne))
	r.POST("dsapi/find", mvc.Returns(i.API.Find))
	r.POST("dsapi/page", mvc.Returns(i.API.Page))
	r.POST("dsapi/create", mvc.Returns(i.API.Create))
	r.POST("dsapi/update", mvc.Returns(i.API.Update))
	r.POST("dsapi/delete", mvc.Returns(i.API.Delete))
}
