package system

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/app/system/controller"
	"lab-api/app/system/service"
	"lab-api/config"
)

var App = fx.Options(service.Provides, controller.Provides, fx.Invoke(Routes))

type Dependency struct {
	fx.In

	Config config.Config
	Authx  *authx.Authx
	Cookie *cookie.Cookie

	*controller.Index
	*controller.Resource
	*controller.Admin
}

func Routes(r *gin.Engine, d Dependency) {
	s := r.Group("/system")
	auth := authMiddleware(d.Authx.Make("system"), d.Cookie)

	s.POST("auth", crud.Bind(d.Index.Login))
	s.GET("auth", crud.Bind(d.Index.Verify))
	s.GET("code", auth, crud.Bind(d.Index.Code))
	s.PUT("auth", auth, crud.Bind(d.Index.RefreshToken))
	s.DELETE("auth", auth, crud.Bind(d.Index.Logout))

	resource := s.Group("resource", auth)
	{
		resource.POST("originLists", crud.Bind(d.Resource.OriginLists))
		resource.POST("lists", crud.Bind(d.Resource.Lists))
		resource.POST("get", crud.Bind(d.Resource.Get))
		resource.POST("add", crud.Bind(d.Resource.Add))
		resource.POST("edit", crud.Bind(d.Resource.Edit))
		resource.POST("delete", crud.Bind(d.Resource.Delete))
	}
}
