package api

import (
	"context"
	"errors"
	"fmt"
	"server/api/index"
	"server/api/sessions"
	"server/api/teams"
	"server/api/users"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
)

var Provides = wire.NewSet(
	index.Provides,
	sessions.Provides,
	teams.Provides,
	users.Provides,
)

type API struct {
	*common.Inject

	Hertz    *server.Hertz
	Index    *index.Controller
	IndexX   *index.Service
	Sessions *sessions.Controller
	Teams    *teams.Controller
	Users    *users.Controller
	UsersX   *users.Service
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	authx := x.Auth()

	x.Hertz.GET("", x.Index.Ping)
	x.Hertz.POST("login", x.Index.Login)
	x.Hertz.GET("verify", x.Index.Verify)
	x.Hertz.POST("logout", authx, x.Index.Logout)
	r := x.Hertz.Group("", authx)

	binds := [][]interface{}{
		{"GET", "user", x.Index.GetUser},
		{"POST", "user/set_password", x.Index.SetUserPassword},
		// Resource API
		{"CRUD", "teams", x.Teams},
		{"GET", "teams/_exists", x.Users.Exists},
		{"GET", "teams/_search", x.Users.Search},
		{"CRUD", "users", x.Users},
		{"GET", "users/_exists", x.Users.Exists},
		{"GET", "users/_search", x.Users.Search},
		{"POST", "users/set_statuses", x.Users.SetStatuses},
	}
	for _, b := range binds {
		if len(b) != 3 {
			continue
		}
		method, resource := b[0].(string), b[1].(string)
		if method != "CRUD" {
			r.Handle(method, resource, b[2].(func(context.Context, *app.RequestContext)))
		} else {
			controller, ok := b[2].(common.Controller)
			if !ok {
				err = errors.New(fmt.Sprintf(`CRUD[%s]: missing method`, resource))
				return
			}
			r.GET(resource, controller.Find)
			r.GET(fmt.Sprintf(`%s/:id`, resource), controller.FindById)
			r.POST(fmt.Sprintf(`%s/create`, resource), controller.Create)
			r.POST(fmt.Sprintf(`%s/update`, resource), controller.Update)
			r.POST(fmt.Sprintf(`%s/delete`, resource), controller.Delete)
		}
	}
	return x.Hertz, nil
}
