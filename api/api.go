package api

import (
	"context"
	"server/api/index"
	"server/api/orgs"
	"server/api/permissions"
	"server/api/roles"
	"server/api/routes"
	"server/api/sessions"
	"server/api/users"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/fx"
)

var Options = fx.Options(
	fx.Provide(func(i Inject) *API { return &API{Inject: &i} }),
	index.Provides,
	sessions.Provides,
	orgs.Provides,
	permissions.Provides,
	roles.Provides,
	routes.Provides,
	users.Provides,
)

type Inject struct {
	fx.In

	Hertz       *server.Hertz
	Index       *index.Controller
	IndexX      *index.Service
	Orgs        *orgs.Controller
	Permissions *permissions.Controller
	Roles       *roles.Controller
	Routes      *routes.Controller
	Users       *users.Controller
	UsersX      *users.Service
}

type API struct{ *Inject }

func (x *API) Initialize(ctx context.Context) (_ *server.Hertz, err error) {

	authx := x.Auth()

	x.Hertz.GET("", x.Index.Ping)

	m := x.Hertz.Group(``, authx)
	{
		// orgs 模块 -> 标准 CRUD 路由
		m.GET("/orgs/:id", x.Orgs.FindById)
		m.GET("/orgs", x.Orgs.Find)
		m.POST("/orgs/create", x.Orgs.Create)
		m.POST("/orgs/update", x.Orgs.Update)
		m.POST("/orgs/delete", x.Orgs.Delete)

		// permissions 模块 -> 标准 CRUD 路由
		m.GET("/permissions/:id", x.Permissions.FindById)
		m.GET("/permissions", x.Permissions.Find)
		m.POST("/permissions/create", x.Permissions.Create)
		m.POST("/permissions/update", x.Permissions.Update)
		m.POST("/permissions/delete", x.Permissions.Delete)

		// roles 模块 -> 标准 CRUD 路由
		m.GET("/roles/:id", x.Roles.FindById)
		m.GET("/roles", x.Roles.Find)
		m.POST("/roles/create", x.Roles.Create)
		m.POST("/roles/update", x.Roles.Update)
		m.POST("/roles/delete", x.Roles.Delete)

		// routes 模块 -> 标准 CRUD 路由
		m.GET("/routes/:id", x.Routes.FindById)
		m.GET("/routes", x.Routes.Find)
		m.POST("/routes/create", x.Routes.Create)
		m.POST("/routes/update", x.Routes.Update)
		m.POST("/routes/delete", x.Routes.Delete)

		// users 模块 -> 标准 CRUD 路由
		m.GET("/users/:id", x.Users.FindById)
		m.GET("/users", x.Users.Find)
		m.POST("/users/create", x.Users.Create)
		m.POST("/users/update", x.Users.Update)
		m.POST("/users/delete", x.Users.Delete)
	}

	return x.Hertz, nil
}
