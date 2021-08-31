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

	resource := s.Group("resource", auth)
	{
		resource.POST("originLists", mvc.Bind(d.Resource.OriginLists))
		resource.POST("get", mvc.Bind(d.Resource.Get))
		resource.POST("add", mvc.Bind(d.Resource.Add))
		resource.POST("edit", mvc.Bind(d.Resource.Edit))
		resource.POST("delete", mvc.Bind(d.Resource.Delete))
	}
	role := s.Group("role", auth)
	{
		role.POST("originLists", mvc.Bind(d.Role.OriginLists))
		role.POST("lists", mvc.Bind(d.Role.Lists))
		role.POST("get", mvc.Bind(d.Role.Get))
		role.POST("add", mvc.Bind(d.Role.Add))
		role.POST("edit", mvc.Bind(d.Role.Edit))
		role.POST("delete", mvc.Bind(d.Role.Delete))
	}
	admin := s.Group("admin", auth)
	{
		admin.POST("originLists", mvc.Bind(d.Admin.OriginLists))
		admin.POST("lists", mvc.Bind(d.Admin.Lists))
		admin.POST("get", mvc.Bind(d.Admin.Get))
		admin.POST("add", mvc.Bind(d.Admin.Add))
		admin.POST("edit", mvc.Bind(d.Admin.Edit))
		admin.POST("delete", mvc.Bind(d.Admin.Delete))
	}
	return &Routes{}
}
