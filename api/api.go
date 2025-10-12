package api

import (
	"context"
	"errors"
	"fmt"
	"server/api/index"
	"server/api/jobs"
	"server/api/schedulers"
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
	jobs.Provides,
	schedulers.Provides,
	sessions.Provides,
	teams.Provides,
	users.Provides,
)

type API struct {
	*common.Inject

	Hertz      *server.Hertz
	Index      *index.Controller
	IndexX     *index.Service
	Jobs       *jobs.Controller
	Schedulers *schedulers.Controller
	Sessions   *sessions.Controller
	Teams      *teams.Controller
	Users      *users.Controller
	UsersX     *users.Service
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
		{"CRUD", "jobs", x.Jobs},
		{"GET", "jobs/_search", x.Jobs.Search},
		{"POST", "jobs/set_statuses", x.Jobs.SetStatuses},
		{"CRUD", "schedulers", x.Schedulers},
		{"GET", "schedulers/_exists", x.Schedulers.Exists},
		{"GET", "schedulers/_search", x.Schedulers.Search},
		{"POST", "schedulers/set_statuses", x.Schedulers.SetStatuses},
		{"GET", "sessions", x.Sessions.Lists},
		{"POST", "sessions/kick", x.Sessions.Kick},
		{"POST", "sessions/clear", x.Sessions.Lists},
		{"CRUD", "teams", x.Teams},
		{"GET", "teams/_exists", x.Teams.Exists},
		{"GET", "teams/_search", x.Teams.Search},
		{"POST", "teams/set_statuses", x.Teams.SetStatuses},
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
