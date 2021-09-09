package system

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/mvc"
	"go.uber.org/fx"
	"lab-api/app/system/admin"
	"lab-api/app/system/index"
	"lab-api/app/system/resource"
	"lab-api/app/system/role"
	"lab-api/common"
)

var Options = fx.Options(
	index.Provides,
	resource.Provides,
	role.Provides,
	admin.Provides,
	fx.Invoke(Routes),
)

type Inject struct {
	common.App

	Index    *index.Controller
	Resource *resource.Controller
	Role     *role.Controller
	Admin    *admin.Controller
}

func Routes(r *gin.Engine, i Inject) {
	s := r.Group("/system")
	auth := authMiddleware(i.Authx.Make("system"), i.Cookie)

	s.POST("login", mvc.Bind(i.Index.Login))
	s.POST("verify", mvc.Bind(i.Index.Verify))
	s.POST("code", auth, mvc.Bind(i.Index.Code))
	s.POST("refresh", auth, mvc.Bind(i.Index.RefreshToken))
	s.POST("logout", auth, mvc.Bind(i.Index.Logout))
	s.POST("resource", auth, mvc.Bind(i.Index.Resource))

	mvc.Crud(s.Group("resource", auth), i.Resource)
	mvc.Crud(s.Group("role", auth), i.Role)
	mvc.Crud(s.Group("admin", auth), i.Admin)
}
