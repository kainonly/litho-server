package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kainonly/go-bit/authx"
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/mvc"
	"lab-api/app/system/controller"
	"lab-api/app/system/service"
	"lab-api/common"
)

var Provides = wire.NewSet(
	service.Provides,
	controller.Provides,
	wire.Struct(new(Dependency), "*"),
	NewRoutes,
)

type Dependency struct {
	Config common.Config
	Cookie *cookie.Cookie
	Authx  *authx.Authx

	*controller.Index
	*controller.Resource
	*controller.Role
	*controller.Admin
}

type Routes struct{}

func NewRoutes(r *gin.Engine, d *Dependency) *Routes {
	s := r.Group("/system")
	auth := authMiddleware(d.Authx.Make("system"), d.Cookie)

	s.POST("login", mvc.Bind(d.Index.Login))
	s.POST("verify", mvc.Bind(d.Index.Verify))
	s.POST("code", auth, mvc.Bind(d.Index.Code))
	s.POST("refresh", auth, mvc.Bind(d.Index.RefreshToken))
	s.POST("logout", auth, mvc.Bind(d.Index.Logout))
	s.POST("resource", auth, mvc.Bind(d.Index.Resource))

	mvc.Crud(s.Group("resource", auth), d.Resource.Crud)
	mvc.Crud(s.Group("role", auth), d.Role.Crud)
	mvc.Crud(s.Group("admin", auth), d.Admin.Crud)
	return &Routes{}
}
