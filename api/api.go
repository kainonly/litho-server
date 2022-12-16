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
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/projects"
	"github.com/weplanx/server/common"
	"github.com/weplanx/utils/dsl"
	"github.com/weplanx/utils/helper"
	"github.com/weplanx/utils/kv"
	"github.com/weplanx/utils/sessions"
	"net/http"
)

var Provides = wire.NewSet(
	index.Provides,
	projects.Provides,
	kv.Provides,
	sessions.Provides,
	dsl.Provides,
)

type API struct {
	*common.Inject

	Hertz    *server.Hertz
	Index    *index.Controller
	Projects *projects.Controller
	KV       *kv.Controller
	Sessions *sessions.Controller
	DSL      *dsl.Controller
}

func (x *API) Routes(h *server.Hertz) (err error) {
	auth := x.AuthGuard()
	h.GET("", x.Index.Ping)
	h.POST("login", x.Index.Login)
	h.GET("verify", x.Index.Verify)
	h.GET("code", auth, x.Index.GetRefreshCode)
	h.POST("refresh_token", auth, x.Index.RefreshToken)
	h.POST("logout", auth, x.Index.Logout)

	user := h.Group("user", auth)
	{
		user.GET("", x.Index.GetUser)
		user.POST("", x.Index.SetUser)
	}

	//projects := h.Group("projects", auth)
	//{
	//	projects.GET("")
	//	projects.POST("")
	//	projects.PATCH(":id")
	//	projects.DELETE(":id")
	//}

	helper.BindKV(h.Group("values", auth), x.KV)
	helper.BindSessions(h.Group("sessions", auth), x.Sessions)
	helper.BindDSL(h.Group(":collection", auth), x.DSL)
	return
}

// AuthGuard
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

		claims, err := x.Index.IndexService.Verify(ctx, string(ts))
		if err != nil {
			c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
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

// AccessLogs 日志
//func (x *API) AccessLogs() app.HandlerFunc {
//	return func(ctx context.Context, c *app.RequestContext) {
//		start := time.Now()
//		c.Next(ctx)
//		end := time.Now()
//		latency := end.Sub(start).Microseconds
//		x.Transfer.Publish(context.Background(), "access", transfer.Payload{
//			Metadata: map[string]interface{}{
//				"method": string(c.Request.Header.Method()),
//				"host":   string(c.Request.Host()),
//				"path":   string(c.Request.Path()),
//				"status": c.Response.StatusCode(),
//				"ip":     c.ClientIP(),
//			},
//			Data: map[string]interface{}{
//				"user_agent": string(c.Request.Header.UserAgent()),
//				"query":      c.Request.QueryString(),
//				"body":       string(c.Request.Body()),
//				"cost":       latency(),
//			},
//			Timestamp: start,
//		})
//	}
//}

// ErrHandler 错误处理
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

// Initialize 初始化
func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	h = x.Hertz
	//h.Use(x.AccessLogs())
	h.Use(x.ErrHandler())
	// 加载自定义验证
	helper.RegValidate()
	// 订阅动态配置
	updated := make(chan *kv.DynamicValues)
	go x.KV.KVService.Sync(&kv.SyncOption{
		Updated: updated,
	})
	// 传输指标
	//if err = x.Transfer.Set(ctx, transfer.LogOption{
	//	Key:         "access",
	//	Description: "请求日志",
	//	TTL:         15552000,
	//}); err != nil {
	//	return
	//}
	//go func() {
	//	for {
	//		select {
	//		case <-updated:
	//			if err = x.DSL.DSLService.Load(ctx); err != nil {
	//				return
	//			}
	//		}
	//	}
	//}()
	return
}
