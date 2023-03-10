package api

import (
	"context"
	"github.com/bytedance/go-tagexpr/v2/binding"
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/google/wire"
	"github.com/weplanx/server/api/feishu"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/monitor"
	"github.com/weplanx/server/api/projects"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/common"
	"github.com/weplanx/transfer"
	"github.com/weplanx/utils/csrf"
	"github.com/weplanx/utils/dsl"
	"github.com/weplanx/utils/helper"
	"github.com/weplanx/utils/kv"
	"github.com/weplanx/utils/passport"
	"github.com/weplanx/utils/sessions"
	"net/http"
	"os"
	"time"
)

var Provides = wire.NewSet(
	index.Provides,
	kv.Provides,
	sessions.Provides,
	dsl.Provides,
	projects.Provides,
	feishu.Provides,
	tencent.Provides,
	monitor.Provides,
)

type API struct {
	*common.Inject

	Hertz    *server.Hertz
	Csrf     *csrf.Csrf
	KV       *kv.Controller
	Sessions *sessions.Controller
	Transfer *transfer.Transfer
	DSL      *dsl.Controller

	Index    *index.Controller
	Projects *projects.Controller
	Feishu   *feishu.Controller
	Tencent  *tencent.Controller
	Monitor  *monitor.Controller
	MonitorX *monitor.Service
}

func (x *API) Routes(h *server.Hertz) (err error) {
	release := os.Getenv("MODE") == "release"
	csrfToken := x.Csrf.VerifyToken(!release)
	auth := x.AuthGuard()

	h.GET("", x.Index.Ping)
	h.POST("login", csrfToken, x.Index.Login)
	h.GET("verify", x.Index.Verify)
	h.GET("code", auth, x.Index.GetRefreshCode)
	h.POST("refresh_token", csrfToken, auth, x.Index.RefreshToken)
	h.POST("logout", csrfToken, auth, x.Index.Logout)

	_user := h.Group("user", csrfToken, auth)
	{
		_user.GET("", x.Index.GetUser)
		_user.POST("", x.Index.SetUser)
		_user.DELETE(":key", x.Index.UnsetUser)
	}

	h.GET("options", x.Index.Options)

	_feishu := h.Group("feishu")
	{
		_feishu.POST("", x.Feishu.Challenge)
		_feishu.GET("", x.Feishu.OAuth)
		_feishu.POST("tasks", x.Feishu.CreateTasks)
		_feishu.GET("tasks", x.Feishu.GetTasks)
	}

	_tencent := h.Group("tencent", auth)
	{
		_tencent.GET("cos-presigned", x.Tencent.CosPresigned)
		_tencent.GET("cos-image-info", x.Tencent.ImageInfo)
	}

	_monitor := h.Group("monitor", auth)
	{
		_monitor.GET("cgo_calls", x.Monitor.GetCgoCalls)
		_monitor.GET("mongo_uptime", x.Monitor.GetMongoUptime)
		_monitor.GET("mongo_available_connections", x.Monitor.GetMongoAvailableConnections)
		_monitor.GET("mongo_open_connections", x.Monitor.GetMongoOpenConnections)
		_monitor.GET("mongo_commands_per_second", x.Monitor.GetMongoCommandsPerSecond)
		_monitor.GET("mongo_query_operations", x.Monitor.GetMongoQueryOperations)
		_monitor.GET("mongo_document_operations", x.Monitor.GetMongoDocumentOperations)
		_monitor.GET("mongo_flushes", x.Monitor.GetMongoFlushes)
		_monitor.GET("mongo_network_io", x.Monitor.GetMongoNetworkIO)
	}

	helper.BindKV(h.Group("values", csrfToken, auth), x.KV)
	helper.BindSessions(h.Group("sessions", csrfToken, auth), x.Sessions)
	helper.BindDSL(h.Group(":collection", csrfToken, auth), x.DSL)

	return
}

func (x *API) AuthGuard() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := c.Cookie("access_token")
		if ts == nil {
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": "authentication has expired, please log in again",
			})
			return
		}

		claims, err := x.Index.Service.Verify(ctx, string(ts))
		if err != nil {
			c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteLaxMode, true, true)
			c.AbortWithStatusJSON(401, utils.H{
				"code":    0,
				"message": index.MsgAuthenticationExpired,
			})
			return
		}

		c.Set("identity", claims)
		c.Next(ctx)
	}
}

func (x *API) AccessLogs() app.HandlerFunc {
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
		x.Transfer.Publish(context.Background(), "access", transfer.Payload{
			Timestamp: now,
			Metadata: map[string]interface{}{
				"method":    method,
				"path":      string(c.Request.Path()),
				"user_id":   userId,
				"client_ip": c.ClientIP(),
			},
			Data: map[string]interface{}{
				"status":     c.Response.StatusCode(),
				"user_agent": string(c.Request.Header.UserAgent()),
			},
		})
	}
}

func (x *API) ErrHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		if err.IsType(errors.ErrorTypePublic) {
			statusCode := http.StatusBadRequest
			result := utils.H{"message": err.Error()}
			if meta, ok := err.Meta.(map[string]interface{}); ok {
				if meta["statusCode"] != nil {
					statusCode = meta["statusCode"].(int)
				}
				if meta["code"] != nil {
					result["code"] = meta["code"]
				}
			}
			c.JSON(statusCode, result)
			return
		}

		switch e := err.Err.(type) {
		case decoder.SyntaxError:
			c.JSON(http.StatusBadRequest, utils.H{
				"code":    0,
				"message": e.Description(),
			})
			break
		case *binding.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"code":    0,
				"message": e.Error(),
			})
			break
		case *validator.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"code":    0,
				"message": e.Error(),
			})
			break
		default:
			logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
	}
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz

	h.Use(x.AccessLogs())
	h.Use(x.ErrHandler())
	helper.RegValidate()

	updated := make(chan *kv.DynamicValues)
	go x.KV.KVService.Sync(&kv.SyncOption{
		Updated: updated,
	})

	if err = x.Transfer.Set(ctx, transfer.LogOption{
		Key:         "access",
		Description: "Access Log Stream",
	}); err != nil {
		return
	}

	go func() {
		for {
			select {
			case <-updated:
				if err = x.DSL.Service.Load(ctx); err != nil {
					return
				}
			}
		}
	}()
	return
}
