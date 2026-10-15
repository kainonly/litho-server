package api

import (
	"context"
	"fmt"
	"litho-api/api/audits"
	"litho-api/api/index"
	"litho-api/api/login_audits"
	"litho-api/api/orgs"
	"litho-api/api/permissions"
	"litho-api/api/roles"
	"litho-api/api/routes"
	"litho-api/api/sessions"
	"litho-api/api/tencent"
	"litho-api/api/users"
	"litho-api/common"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/goforj/wire"
	"github.com/kainonly/go/csrf"
)

var Provides = wire.NewSet(
	audits.Provides,
	orgs.Provides,
	index.Provides,
	login_audits.Provides,
	permissions.Provides,
	roles.Provides,
	routes.Provides,
	sessions.Provides,
	tencent.Provides,
	users.Provides,
)

type API struct {
	*common.Inject

	Hertz *server.Hertz
	Csrf  *csrf.Csrf

	Index       *index.Controller
	IndexX      *index.Service
	Audits      *audits.Controller
	LoginAudits *login_audits.Controller
	Orgs        *orgs.Controller
	Permissions *permissions.Controller
	Roles       *roles.Controller
	RolesX      *roles.Service
	Routes      *routes.Controller
	Sessions    *sessions.Controller
	Users       *users.Controller
	UsersX      *users.Service
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	_auth, _audit := x.Auth(), x.Audit()

	x.Hertz.GET("", x.Index.Ping)
	x.Hertz.POST("login", x.Index.Login)
	x.Hertz.GET("login/sms", x.Index.GetLoginSms)
	x.Hertz.POST("login/sms", x.Index.LoginSms)
	x.Hertz.GET("verify", x.Index.Verify)
	x.Hertz.POST("logout", _auth, x.Index.Logout)
	r := x.Hertz.Group("", _auth, _audit)

	binds := [][]any{
		{"GET", "layout", x.Index.GetLayout},
		{"GET", "user", x.Index.GetUser},
		{"GET", "user/phone_code", x.Index.GetUserPhoneCode},
		{"GET", "user/history", x.Index.GetHistory},
		{"POST", "user/set", x.Index.SetUser},
		{"POST", "user/set_phone", x.Index.SetUserPhone},
		{"POST", "user/set_password", x.Index.SetUserPassword},
		{"POST", "user/unset", x.Index.UnsetUser},

		{"GET", "audits", x.Audits.Find},
		{"GET", "audits/schema", x.Audits.GetSchema},
		{"GET", "login_audits", x.LoginAudits.Find},
		{"GET", "orgs/_search", x.Orgs.Search},
		{"CRUD", "orgs", x.Orgs},
		{"GET", "permissions/_search", x.Permissions.Search},
		{"CRUD", "permissions", x.Permissions},
		{"GET", "roles/_super", x.Roles.GetSuper},
		{"GET", "roles/_search", x.Roles.Search},
		{"CRUD", "roles", x.Roles},
		{"POST", "roles/set_strategy", x.Roles.SetStrategy},
		{"POST", "roles/sort", x.Roles.Sort},
		{"GET", "routes/_resources", x.Routes.GetResources},
		{"GET", "routes/_search", x.Routes.Search},
		{"CRUD", "routes", x.Routes},
		{"POST", "routes/regroup", x.Routes.Regroup},
		{"POST", "routes/set_strategy", x.Routes.SetStrategy},
		{"GET", "sessions", x.Sessions.Lists},
		{"POST", "sessions/kick", x.Sessions.Kick},
		{"POST", "sessions/clear", x.Sessions.Clear},
		{"GET", "users/_exists", x.Users.Exists},
		{"GET", "users/_search", x.Users.Search},
		{"CRUD", "users", x.Users},
		{"POST", "users/set_statuses", x.Users.SetStatuses},
		{"POST", "users/set_roles", x.Users.SetRoles},
		{"POST", "users/set_orgs", x.Users.SetOrgs},
	}

	writePaths := make([]string, 0)
	for _, b := range binds {
		method, path := b[0].(string), b[1].(string)
		if method == "CRUD" {
			controller, ok := b[2].(common.Controller)
			if !ok {
				return nil, fmt.Errorf("CRUD[%s]: missing method", path)
			}
			r.GET(path, controller.Find)
			r.GET(path+"/:id", controller.FindById)
			r.POST(path+"/create", controller.Create)
			r.POST(path+"/update", controller.Update)
			r.POST(path+"/delete", controller.Delete)
			writePaths = append(writePaths,
				path+"/create", path+"/update", path+"/delete",
			)
			continue
		}
		handler, ok := b[2].(func(context.Context, *app.RequestContext))
		if !ok {
			return nil, fmt.Errorf("%s %s: invalid handler", method, path)
		}
		r.Handle(method, path, handler)
		if method == "POST" && !x.IgnoreRoutes(strings.Split(path, "/")) {
			writePaths = append(writePaths, path)
		}
	}
	if err = x.V.Resources.ValidateRoutes(writePaths); err != nil {
		return nil, err
	}
	return x.Hertz, nil
}
