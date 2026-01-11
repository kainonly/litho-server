package api

import (
	"context"
	"server/api/index"
	"server/api/orgs"
	"server/api/roles"
	"server/api/users"

	"github.com/cloudwego/hertz/pkg/app/server"
	"go.uber.org/fx"
)

var Options = fx.Options(
	index.Provides,
	orgs.Provides,
	roles.Provides,
	users.Provides,
	fx.Invoke(Routes),
)

func Routes(lc fx.Lifecycle, h *server.Hertz, index *index.Controller, orgs *orgs.Controller, roles *roles.Controller, users *users.Controller) {
	h.GET("", index.Ping)

	// orgs 模块 -> 标准 CRUD 路由
	h.GET("/orgs/:id", orgs.FindById)
	h.GET("/orgs", orgs.Find)
	h.POST("/orgs/create", orgs.Create)
	h.POST("/orgs/update", orgs.Update)
	h.POST("/orgs/delete", orgs.Delete)

	// roles 模块 -> 标准 CRUD 路由
	h.GET("/roles/:id", roles.FindById)
	h.GET("/roles", roles.Find)
	h.POST("/roles/create", roles.Create)
	h.POST("/roles/update", roles.Update)
	h.POST("/roles/delete", roles.Delete)

	// users 模块 -> 标准 CRUD 路由
	h.GET("/users/:id", users.FindById)
	h.GET("/users", users.Find)
	h.POST("/users/create", users.Create)
	h.POST("/users/update", users.Update)
	h.POST("/users/delete", users.Delete)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go h.Spin()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if h != nil {
				return h.Shutdown(ctx)
			}
			return nil
		},
	})
}
