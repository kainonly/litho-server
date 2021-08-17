package system

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/mvc"
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

	s.POST("login", mvc.Bind(d.Index.Login))
	s.POST("verify", mvc.Bind(d.Index.Verify))
	s.POST("code", auth, mvc.Bind(d.Index.Code))
	s.POST("refresh", auth, mvc.Bind(d.Index.RefreshToken))
	s.POST("logout", auth, mvc.Bind(d.Index.Logout))

	resource := s.Group("resource", auth)
	{
		resource.POST("originLists", mvc.Bind(d.Resource.OriginLists))
		resource.POST("get", mvc.Bind(d.Resource.Get))
		resource.POST("add", mvc.Bind(d.Resource.Add))
		resource.POST("edit", mvc.Bind(d.Resource.Edit))
		resource.POST("delete", mvc.Bind(d.Resource.Delete))
	}
}
