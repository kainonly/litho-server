package system

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/app/system/controller"
	"lab-api/app/system/service"
	"lab-api/common"
)

var App = fx.Options(service.Provides, controller.Provides, fx.Invoke(Routes))

type Dependency struct {
	fx.In

	App    *common.App
	Authx  *authx.Authx
	Cookie *cookie.Cookie

	*controller.Index
	*controller.Resource
}

func Routes(r *gin.Engine, d Dependency) {
	s := r.Group("/system")
	auth := authMiddleware(d.Authx.Make("system"), d.Cookie)

	s.POST("login", mvc.Bind(d.Index.Login))
	s.POST("verify", mvc.Bind(d.Index.Verify))
	s.POST("code", auth, mvc.Bind(d.Index.Code))
	s.POST("refresh", auth, mvc.Bind(d.Index.RefreshToken))
	s.POST("logout", auth, mvc.Bind(d.Index.Logout))
	s.POST("resource", auth, mvc.Bind(d.Index.Resource))

	mvc.Crud(s.Group("resource", auth), d.Resource.Crud)
}
