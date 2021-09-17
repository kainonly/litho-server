package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"lab-api/app/system/index"
	"lab-api/common"
)

var Options = fx.Options(
	index.Provides,
	//resource.Provides,
	//role.Provides,
	//admin.Provides,
	fx.Invoke(Routes),
)

type Inject struct {
	common.App

	Index *index.Controller
	//Resource *resource.Controller
	//Role     *role.Controller
	//Admin    *admin.Controller
}

func Routes(r *gin.Engine, i Inject) {
	//s := r.Group("system")
	//if os.Getenv("GIN_MODE") != "release" {
	//	devtoolsRoute := s.Group("devtools")
	//	devtoolsRoute.POST("setup", mvc.Returns(i.DevTools.Setup))
	//	devtoolsRoute.POST("sync", mvc.Returns(i.DevTools.Sync))
	//	devtoolsRoute.POST("migrate", mvc.Returns(i.DevTools.Migrate))
	//}
	//auth := authMiddleware(i.Authx.Make("system"), i.Cookie)
	//s.POST("login", mvc.Returns(i.Index.Login))
	//s.POST("verify", mvc.Returns(i.Index.Verify))
	//s.POST("code", auth, mvc.Returns(i.Index.Code))
	//s.POST("refresh", auth, mvc.Returns(i.Index.RefreshToken))
	//s.POST("logout", auth, mvc.Returns(i.Index.Logout))
	//s.POST("resource", auth, mvc.Returns(i.Index.Resource))
	//resourceRoute := s.Group("resource")
	//{
	//	mvc.Crud(resourceRoute, i.Resource)
	//}
	//roleRoute := s.Group("role", auth)
	//{
	//	mvc.Crud(roleRoute, i.Role)
	//}
	//adminRoute := s.Group("admin", auth)
	//{
	//	mvc.Crud(adminRoute, i.Admin)
	//}
}
