package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/kainonly/go-bit/mvc"
	"lab-api/app/system/index"
	"lab-api/app/system/resource"
	"lab-api/common"
)

var Provides = wire.NewSet(
	wire.Struct(new(Inject), "*"),
	index.Provides,
	resource.Provides,
	NewRoutes,
)

type Inject struct {
	Index    *index.Controller
	Resource *resource.Controller
}

type Routes struct{}

func NewRoutes(r *gin.Engine, d *common.Dependency, i *Inject) *Routes {
	s := r.Group("/system")
	auth := authMiddleware(d.Authx.Make("system"), d.Cookie)

	s.POST("login", mvc.Bind(i.Index.Login))
	s.POST("verify", mvc.Bind(i.Index.Verify))
	s.POST("code", auth, mvc.Bind(i.Index.Code))
	s.POST("refresh", auth, mvc.Bind(i.Index.RefreshToken))
	s.POST("logout", auth, mvc.Bind(i.Index.Logout))
	s.POST("resource", auth, mvc.Bind(i.Index.Resource))

	mvc.Crud(s.Group("resource", auth), i.Resource.Crud)
	return &Routes{}
}
