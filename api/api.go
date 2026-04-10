package api

import (
	"context"
	"server/api/departments"
	"server/api/index"
	"server/api/orders"
	"server/api/permissions"
	"server/api/products"
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
	permissions.Provides,
	index.Provides,
	departments.Provides,
	orders.Provides,
	products.Provides,
	resources.Provides,
	roles.Provides,
	routes.Provides,
	sessions.Provides,
	users.Provides,
)

type API struct {
	*common.Inject

	Hertz       *server.Hertz
	Csrf        *csrf.Csrf
	Permissions *permissions.Controller
	Index       *index.Controller
	IndexX      *index.Service
	Departments *departments.Controller
	Orders      *orders.Controller
	Products    *products.Controller
	Resources   *resources.Controller
	Roles       *roles.Controller
	RolesX      *roles.Service
	Routes      *routes.Controller
	Sessions    *sessions.Controller
	Users       *users.Controller
	UsersX      *users.Service
}

func (x *API) Initialize(ctx context.Context) (_ *server.Hertz, err error) {
	_auth := x.Auth()

	x.Hertz.GET("", x.Index.Ping)
	x.Hertz.POST("login", x.Index.Login)
	x.Hertz.GET("verify", x.Index.Verify)
	x.Hertz.POST("logout", _auth, x.Index.Logout)
	x.Hertz.GET("layout", _auth, x.Index.GetLayout)
	x.Hertz.GET("user", _auth, x.Index.GetUser)
	x.Hertz.POST("set_user", _auth, x.Index.SetUser)
	x.Hertz.POST("set_user_password", _auth, x.Index.SetUserPassword)
	x.Hertz.POST("set_user_phone", _auth, x.Index.SetUserPhone)
	x.Hertz.POST("unset_user", _auth, x.Index.UnsetUser)

	m := x.Hertz.Group(``, _auth)
	{
		// permissions 模块 -> 只读，由 sync-meta 维护
		m.GET("/permissions/:id", x.Permissions.FindById)
		m.GET("/permissions", x.Permissions.Find)

		// departments 模块 -> 标准 CRUD 路由
		m.GET("/departments/:id", x.Departments.FindById)
		m.GET("/departments", x.Departments.Find)
		m.POST("/departments/create", x.Departments.Create)
		m.POST("/departments/update", x.Departments.Update)
		m.POST("/departments/delete", x.Departments.Delete)

		// resources 模块 -> 只读，由 sync-meta 维护
		m.GET("/resources/:id", x.Resources.FindById)
		m.GET("/resources", x.Resources.Find)
		m.GET("/resources/_search", x.Resources.Search)

		// roles 模块 -> 标准 CRUD 路由
		m.GET("/roles/:id", x.Roles.FindById)
		m.GET("/roles", x.Roles.Find)
		m.POST("/roles/create", x.Roles.Create)
		m.POST("/roles/update", x.Roles.Update)
		m.POST("/roles/delete", x.Roles.Delete)
		m.POST("/roles/sort", x.Roles.Sort)
		m.POST("/roles/set_strategy", x.Roles.SetStrategy)

		// routes 模块 -> 标准 CRUD 路由
		m.GET("/routes/:id", x.Routes.FindById)
		m.GET("/routes", x.Routes.Find)
		m.GET("/routes/_search", x.Routes.Search)
		m.POST("/routes/create", x.Routes.Create)
		m.POST("/routes/update", x.Routes.Update)
		m.POST("/routes/delete", x.Routes.Delete)
		m.POST("/routes/sort", x.Routes.Sort)
		m.POST("/routes/regroup", x.Routes.Regroup)

		// sessions 模块 -> 会话管理路由
		m.GET("/sessions", x.Sessions.Lists)
		m.POST("/sessions/kick", x.Sessions.Kick)
		m.POST("/sessions/clear", x.Sessions.Clear)

		// users 模块 -> 标准 CRUD 路由
		m.GET("/users/:id", x.Users.FindById)
		m.GET("/users", x.Users.Find)
		m.GET("/users/_exists", x.Users.Exists)
		m.GET("/users/_search", x.Users.Search)
		m.POST("/users/create", x.Users.Create)
		m.POST("/users/update", x.Users.Update)
		m.POST("/users/delete", x.Users.Delete)
		m.POST("/users/set_roles", x.Users.SetRoles)
		m.POST("/users/set_statuses", x.Users.SetStatuses)

		// products 模块 -> 标准 CRUD 路由
		m.GET("/products/:id", x.Products.FindById)
		m.GET("/products", x.Products.Find)
		m.GET("/products/_search", x.Products.Search)
		m.POST("/products/create", x.Products.Create)
		m.POST("/products/update", x.Products.Update)
		m.POST("/products/delete", x.Products.Delete)

		// orders 模块 -> 标准 CRUD 路由
		m.GET("/orders/:id", x.Orders.FindById)
		m.GET("/orders", x.Orders.Find)
		m.POST("/orders/create", x.Orders.Create)
		m.POST("/orders/update", x.Orders.Update)
		m.POST("/orders/delete", x.Orders.Delete)
	}

	return x.Hertz, nil
}
