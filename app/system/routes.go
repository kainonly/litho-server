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
	"log"
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
	//auth := authMiddleware(i.Authx.Make("system"), i.Cookie)
	s.POST("/:model/r/find/one", func(c *gin.Context) {
		var uri struct {
			Model string `uri:"model"`
		}
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		log.Println(uri)
	})
	//s.POST("login", mvc.Bind(i.Index.Login))
	//s.POST("verify", mvc.Bind(i.Index.Verify))
	//s.POST("code", auth, mvc.Bind(i.Index.Code))
	//s.POST("refresh", auth, mvc.Bind(i.Index.RefreshToken))
	//s.POST("logout", auth, mvc.Bind(i.Index.Logout))
	//s.POST("resource", auth, mvc.Bind(i.Index.Resource))

	//mvc.Crud(s.Group("resource"), i.Resource)
	mvc.Crud(s.Group("role"), i.Role)
	mvc.Crud(s.Group("admin"), i.Admin)
}
