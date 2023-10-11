package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	transfer "github.com/weplanx/collector/client"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/go/values"
	"github.com/weplanx/server/api/clusters"
	"github.com/weplanx/server/api/datasets"
	"github.com/weplanx/server/api/imessages"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/lark"
	"github.com/weplanx/server/api/observability"
	"github.com/weplanx/server/api/queues"
	"github.com/weplanx/server/api/schedules"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/api/workflows"
	"github.com/weplanx/server/common"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(values.Controller), "*"),
	wire.Struct(new(sessions.Controller), "*"),
	wire.Struct(new(rest.Controller), "*"),
	index.Provides,
	tencent.Provides,
	lark.Provides,
	clusters.Provides,
	schedules.Provides,
	workflows.Provides,
	queues.Provides,
	imessages.Provides,
	datasets.Provides,
	observability.Provides,
)

type API struct {
	*common.Inject

	Hertz                *server.Hertz
	Csrf                 *csrf.Csrf
	Values               *values.Controller
	Sessions             *sessions.Controller
	Rest                 *rest.Controller
	Index                *index.Controller
	IndexService         *index.Service
	Tencent              *tencent.Controller
	TencentService       *tencent.Service
	Lark                 *lark.Controller
	LarkService          *lark.Service
	Clusters             *clusters.Controller
	ClustersService      *clusters.Service
	Schedules            *schedules.Controller
	SchedulesService     *schedules.Service
	Workflows            *workflows.Controller
	WorkflowsService     *workflows.Service
	Queues               *queues.Controller
	QueuesService        *queues.Service
	Imessages            *imessages.Controller
	ImessagesService     *imessages.Service
	Datasets             *datasets.Controller
	DatasetsService      *datasets.Service
	Observability        *observability.Controller
	ObservabilityService *observability.Service
}

func (x *API) Routes(h *server.Hertz) (err error) {
	csrfToken := x.Csrf.VerifyToken(!x.V.IsRelease())
	auth := x.AuthGuard()
	audit := x.Audit()

	h.GET("", x.Index.Ping)
	_login := h.Group("login", csrfToken)
	{
		_login.POST("", x.Index.Login)
		_login.GET("sms", x.Index.GetLoginSms)
		_login.POST("sms", x.Index.LoginSms)
		_login.POST("totp", x.Index.LoginTotp)
	}
	h.GET("forget_code", x.Index.GetForgetCode)
	h.POST("forget_reset", csrfToken, x.Index.ForgetReset)
	h.GET("verify", csrfToken, x.Index.Verify)
	h.GET("refresh_code", csrfToken, auth, x.Index.GetRefreshCode)
	h.POST("refresh_token", csrfToken, auth, x.Index.RefreshToken)
	h.POST("logout", csrfToken, auth, x.Index.Logout)
	h.GET("options", x.Index.Options)

	m := []app.HandlerFunc{csrfToken, auth, audit}
	_user := h.Group("user", m...)
	{
		_user.GET("", x.Index.GetUser)
		_user.PATCH("", x.Index.SetUser)
		_user.POST("password", x.Index.SetUserPassword)
		_user.GET("phone_code", x.Index.GetUserPhoneCode)
		_user.POST("phone", x.Index.SetUserPhone)
		_user.GET("totp", x.Index.GetUserTotp)
		_user.POST("totp", x.Index.SetUserTotp)
		_user.DELETE(":key", x.Index.UnsetUser)
	}
	_values := h.Group("values", m...)
	{
		_values.GET("", x.Values.Get)
		_values.PATCH("", x.Values.Set)
		_values.DELETE(":key", x.Values.Remove)
	}
	_sessions := h.Group("sessions", m...)
	{
		_sessions.GET("", x.Sessions.Lists)
		_sessions.DELETE(":uid", x.Sessions.Remove)
		_sessions.POST("clear", x.Sessions.Clear)
	}
	_db := h.Group("db", csrfToken, auth)
	{
		_db.GET(":collection/:id", x.Rest.FindById)
		_db.POST(":collection/create", audit, x.Rest.Create)
		_db.POST(":collection/bulk_create", audit, x.Rest.BulkCreate)
		_db.POST(":collection/size", x.Rest.Size)
		_db.POST(":collection/find", x.Rest.Find)
		_db.POST(":collection/find_one", x.Rest.FindOne)
		_db.POST(":collection/update", audit, x.Rest.Update)
		_db.POST(":collection/bulk_delete", audit, x.Rest.BulkDelete)
		_db.POST(":collection/sort", audit, x.Rest.Sort)
		_db.PATCH(":collection/:id", audit, x.Rest.UpdateById)
		_db.PUT(":collection/:id", audit, x.Rest.Replace)
		_db.DELETE(":collection/:id", audit, x.Rest.Delete)
		_db.POST("transaction", audit, x.Rest.Transaction)
		_db.POST("commit", audit, x.Rest.Commit)
	}
	_tencent := h.Group("tencent", m...)
	{
		_tencent.GET("cos_presigned", x.Tencent.CosPresigned)
		_tencent.GET("cos_image_info", x.Tencent.CosImageInfo)
	}
	h.POST("lark", x.Lark.Challenge)
	h.GET("lark", x.Lark.OAuth)
	_lark := h.Group("lark", m...)
	{
		_lark.POST("tasks", x.Lark.CreateTasks)
		_lark.GET("tasks", x.Lark.GetTasks)
	}
	_clusters := h.Group("clusters", m...)
	{
		_clusters.GET(":id/info", x.Clusters.GetInfo)
		_clusters.GET(":id/nodes", x.Clusters.GetNodes)
	}
	_schedules := h.Group("schedules", m...)
	{
		_schedules.GET(":id/keys", x.Schedules.Keys)
		_schedules.POST("ping", x.Schedules.Ping)
		_schedules.POST("revoke", x.Schedules.Revoke)
		_schedules.POST("state", x.Schedules.State)
	}
	_workflows := h.Group("workflows", m...)
	{
		_workflows.POST("sync", x.Workflows.Sync)
	}
	_queues := h.Group("queues", m...)
	{
		_queues.POST("sync", x.Queues.Sync)
		_queues.POST("destroy", x.Queues.Destroy)
		_queues.GET(":id/info", x.Queues.Info)
		_queues.POST("publish", x.Queues.Publish)
	}
	_imessages := h.Group("imessages", m...)
	{
		_imessages.GET("nodes", x.Imessages.GetNodes)
		_imessages.GET(":id/metrics", x.Imessages.GetMetrics)
		_imessages.PUT(":id/metrics", x.Imessages.CreateMetrics)
		_imessages.DELETE(":id/metrics", x.Imessages.DeleteMetrics)
		_imessages.POST("publish", x.Imessages.Publish)
	}
	_datasets := h.Group("datasets", m...)
	{
		_datasets.GET("", x.Datasets.Lists)
		_datasets.POST("create", x.Datasets.Create)
		_datasets.DELETE(":name", x.Datasets.Delete)
	}
	_observability := h.Group("observability", m...)
	{
		_observability.GET(":name", x.Observability.Exporters)
	}
	return
}

func (x *API) AuthGuard() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": "authentication has expired, please log in again",
			})
			return
		}

		claims, err := x.IndexService.Verify(ctx, string(ts))
		if err != nil {
			common.ClearAccessToken(c)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": common.ErrAuthenticationExpired.Error(),
			})
			return
		}

		c.Set("identity", claims)
		c.Next(ctx)
	}
}

func (x *API) Audit() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		now := time.Now()
		c.Next(ctx)
		method := string(c.Request.Header.Method())
		if method == "GET" {
			return
		}
		var userId string
		if value, ok := c.Get("identity"); ok {
			claims := value.(passport.Claims)
			userId = claims.UserId
		}

		format := map[string]interface{}{
			"body": "json",
		}
		if userId != "" {
			format["metadata.user_id"] = "oid"
		}
		transferCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		x.Transfer.Publish(transferCtx, "logset_operates", transfer.Payload{
			Timestamp: now,
			Data: map[string]interface{}{
				"metadata": map[string]interface{}{
					"method":    method,
					"path":      string(c.Request.Path()),
					"user_id":   userId,
					"client_ip": c.ClientIP(),
				},
				"params": string(c.Request.QueryString()),
				"body":   c.Request.Body(),
				"status": c.Response.StatusCode(),
			},
			XData: format,
		})
	}
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz

	update := make(chan interface{})
	go x.Values.Service.Sync(x.V.Extra, update)
	go x.ValuesChange(update)

	if err = x.Transfer.Set(ctx, transfer.StreamOption{
		Key: "logset_operates",
	}); err != nil {
		return
	}

	go func() {
		if err = x.WorkflowsService.Event(); err != nil {
			hlog.Error(err)
		}
		if err = x.QueuesService.Event(); err != nil {
			hlog.Error(err)
		}
		if err = x.ImessagesService.Event(); err != nil {
			hlog.Error(err)
		}
	}()

	return
}

func (x *API) ValuesChange(ok chan interface{}) {
	for range ok {
		for k, v := range x.V.RestControls {
			if v.Event {
				if _, err := x.JetStream.AddStream(&nats.StreamConfig{
					Name:      x.V.Name("events", k),
					Subjects:  []string{x.V.NameX(".", "events", k)},
					Retention: nats.WorkQueuePolicy,
				}); err != nil {
					hlog.Error(err)
				}
			}
		}
	}
	return
}
