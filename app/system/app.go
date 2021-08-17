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
	*controller.Acl
	*controller.Resource
	*controller.Role
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
	s.POST("resource", auth, mvc.Bind(d.Index.Resource))

	acl := s.Group("acl", auth)
	{
		acl.POST("lists", mvc.Bind(d.Acl.Lists))
		acl.POST("get", mvc.Bind(d.Acl.Get))
		acl.POST("add", mvc.Bind(d.Acl.Add))
		acl.POST("edit", mvc.Bind(d.Acl.Edit))
		acl.POST("delete", mvc.Bind(d.Acl.Delete))
	}
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
}
