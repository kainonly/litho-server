package api

import (
	"context"
	"server/api/caps"
	"server/api/index"
	"server/api/orgs"
	"server/api/resources"
	"server/api/roles"
	"server/api/routes"
	"server/api/sessions"
	"server/api/users"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/goforj/wire"
	"github.com/kainonly/go/csrf"
)

var Provides = wire.NewSet(
	caps.Provides,
	index.Provides,
	orgs.Provides,
	resources.Provides,
	roles.Provides,
	routes.Provides,
	sessions.Provides,
	users.Provides,
)

type API struct {
	*common.Inject

	Hertz     *server.Hertz
	Csrf      *csrf.Csrf
	Caps      *caps.Controller
	Index     *index.Controller
	IndexX    *index.Service
	Orgs      *orgs.Controller
	Resources *resources.Controller
	Roles     *roles.Controller
	Routes    *routes.Controller
	Sessions  *sessions.Controller
	Users     *users.Controller
	UsersX    *users.Service
}

func (x *API) Initialize(ctx context.Context) (_ *server.Hertz, err error) {
	_auth := x.Auth()

	x.Hertz.GET("", x.Index.Ping)
	x.Hertz.POST("login", x.Index.Login)
	x.Hertz.GET("verify", x.Index.Verify)
	x.Hertz.POST("logout", _auth, x.Index.Logout)

	m := x.Hertz.Group(``, _auth)
	{
		// caps 模块 -> 标准 CRUD 路由
		m.GET("/caps/:id", x.Caps.FindById)
		m.GET("/caps", x.Caps.Find)
		m.POST("/caps/create", x.Caps.Create)
		m.POST("/caps/update", x.Caps.Update)
		m.POST("/caps/delete", x.Caps.Delete)

		// orgs 模块 -> 标准 CRUD 路由
		m.GET("/orgs/:id", x.Orgs.FindById)
		m.GET("/orgs", x.Orgs.Find)
		m.POST("/orgs/create", x.Orgs.Create)
		m.POST("/orgs/update", x.Orgs.Update)
		m.POST("/orgs/delete", x.Orgs.Delete)

		// resources 模块 -> 标准 CRUD 路由
		m.GET("/resources/:id", x.Resources.FindById)
		m.GET("/resources", x.Resources.Find)
		m.GET("/resources/_search", x.Resources.Search)
		m.POST("/resources/create", x.Resources.Create)
		m.POST("/resources/update", x.Resources.Update)
		m.POST("/resources/delete", x.Resources.Delete)

		// roles 模块 -> 标准 CRUD 路由
		m.GET("/roles/:id", x.Roles.FindById)
		m.GET("/roles", x.Roles.Find)
		m.POST("/roles/create", x.Roles.Create)
		m.POST("/roles/update", x.Roles.Update)
		m.POST("/roles/delete", x.Roles.Delete)
		m.POST("/roles/sort", x.Roles.Sort)

		// routes 模块 -> 标准 CRUD 路由
		m.GET("/routes/:id", x.Routes.FindById)
		m.GET("/routes", x.Routes.Find)
		m.GET("/routes/_search", x.Routes.Search)
		m.POST("/routes/create", x.Routes.Create)
		m.POST("/routes/update", x.Routes.Update)
		m.POST("/routes/delete", x.Routes.Delete)
		m.POST("/routes/sort", x.Routes.Sort)
		m.POST("/routes/regroup", x.Routes.Regroup)

		// users 模块 -> 标准 CRUD 路由
		m.GET("/users/:id", x.Users.FindById)
		m.GET("/users", x.Users.Find)
		m.GET("/users/_exists", x.Users.Exists)
		m.GET("/users/_search", x.Users.Search)
		m.POST("/users/create", x.Users.Create)
		m.POST("/users/update", x.Users.Update)
		m.POST("/users/delete", x.Users.Delete)
		m.POST("/users/set_roles", x.Users.SetRoles)
		m.POST("/users/set_actives", x.Users.SetActives)
	}

	return x.Hertz, nil
}
