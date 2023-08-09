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
	"github.com/google/wire"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/rest"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/go/values"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/lark"
	"github.com/weplanx/server/api/observability"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/common"
	"github.com/weplanx/transfer"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var Provides = wire.NewSet(
	index.Provides,
	tencent.Provides,
	lark.Provides,
	observability.Provides,
	wire.Struct(new(values.Controller), "*"),
	wire.Struct(new(sessions.Controller), "*"),
	wire.Struct(new(rest.Controller), "*"),
)

type API struct {
	*common.Inject

	Hertz                *server.Hertz
	Csrf                 *csrf.Csrf
	Transfer             *transfer.Transfer
	Values               *values.Controller
	Sessions             *sessions.Controller
	Rest                 *rest.Controller
	Index                *index.Controller
	IndexService         *index.Service
	Tencent              *tencent.Controller
	TencentSerice        *tencent.Service
	Lark                 *lark.Controller
	LarkService          *lark.Service
	Observability        *observability.Controller
	ObservabilityService *observability.Service
}

func (x *API) Routes(h *server.Hertz) (err error) {
	//csrfToken := x.Csrf.VerifyToken(!x.V.IsRelease())
	auth := x.AuthGuard()

	h.GET("", x.Index.Ping)
	h.POST("login", x.Index.Login)
	h.POST("login/totp", x.Index.LoginTotp)
	h.GET("verify", x.Index.Verify)
	h.GET("options", x.Index.Options)
	h.GET("code", auth, x.Index.GetRefreshCode)
	h.POST("refresh_token", auth, x.Index.RefreshToken)
	h.POST("logout", auth, x.Index.Logout)

	m := []app.HandlerFunc{auth, x.Audit()}
	u := h.Group("", m...)
	{
		help.ValuesRoutes(u, x.Values)
		help.SessionsRoutes(u, x.Sessions)
		help.RestRoutes(u.Group("db"), x.Rest)
	}
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
	_observability := h.Group("observability", m...)
	{
		_observability.GET(":name", x.Observability.Exporters)
	}
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
		transferCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
				"params":     string(c.Request.QueryString()),
				"body":       c.Request.Body(),
				"status":     c.Response.StatusCode(),
				"user_agent": string(c.Request.Header.UserAgent()),
			},
			Format: format,
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
		case mongo.ServerError:
			c.JSON(http.StatusInternalServerError, utils.H{
				"code":    0,
				"message": e.Error(),
			})
			break
		default:
			if !x.V.IsRelease() {
				c.JSON(http.StatusInternalServerError, utils.H{
					"code":    0,
					"message": e.Error(),
				})
				break
			}
			logger.Error(err)
			c.Status(http.StatusInternalServerError)
		}
	}
}

func (x *API) Initialize(ctx context.Context) (h *server.Hertz, err error) {
	help.RegValidate()

	h = x.Hertz
	h.Use(x.ErrHandler())

	go x.Values.Service.Sync(x.V.Extra, nil)

	if err = x.Transfer.Set(ctx, transfer.LogOption{
		Key: "logset_operates",
	}); err != nil {
		return
	}

	return
}
