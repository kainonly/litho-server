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
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/hertz-contrib/jwt"
	"github.com/weplanx/server/api/departments"
	"github.com/weplanx/server/api/dsl"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/pages"
	"github.com/weplanx/server/api/roles"
	"github.com/weplanx/server/api/sessions"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/api/values"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/utils/validation"
	"github.com/weplanx/transfer"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

var Provides = wire.NewSet(
	index.Provides,
	values.Provides,
	sessions.Provides,
	dsl.Provides,
	pages.Provides,
	users.Provides,
	roles.Provides,
	departments.Provides,
)

type API struct {
	Values   *common.Values
	Db       *mongo.Database
	Redis    *redis.Client
	Transfer *transfer.Transfer
	Hertz    *server.Hertz

	IndexController   *index.Controller
	IndexService      *index.Service
	ValuesController  *values.Controller
	ValuesService     *values.Service
	SessionController *sessions.Controller
	SessionService    *sessions.Service
	DslController     *dsl.Controller
	DslService        *dsl.Service
	PagesController   *pages.Controller
	PagesService      *pages.Service
	UsersController   *users.Controller
	UsersService      *users.Service
}

func (x *API) Run() (h *server.Hertz, err error) {
	h = x.Hertz
	h.Use(x.AccessLog())
	h.Use(x.ErrHandler())
	return
}

func (x *API) Routes(h *server.Hertz) (auth *jwt.HertzJWTMiddleware, err error) {
	if auth, err = x.Auth(); err != nil {
		return
	}

	h.GET("", x.IndexController.Index)
	h.POST("login", auth.LoginHandler)
	h.GET("code", auth.MiddlewareFunc(), x.IndexController.GetRefreshCode)
	h.POST("refresh_token", auth.MiddlewareFunc(), x.IndexController.VerifyRefreshCode, auth.RefreshHandler)

	h.GET("navs", auth.MiddlewareFunc(), x.IndexController.GetNavs)
	h.GET("options", auth.MiddlewareFunc(), x.IndexController.GetOptions)

	h.GET("user", auth.MiddlewareFunc(), x.IndexController.GetUser)
	h.PATCH("user", auth.MiddlewareFunc(), x.IndexController.SetUser)
	h.POST("unset-user", auth.MiddlewareFunc(), x.IndexController.UnsetUser)
	h.DELETE("user", auth.MiddlewareFunc(), auth.LogoutHandler)

	_values := h.Group("values", auth.MiddlewareFunc())
	{
		_values.GET("", x.ValuesController.Get)
		_values.PATCH("", x.ValuesController.Set)
		_values.DELETE(":key", x.ValuesController.Remove)
	}

	_sessions := h.Group("sessions", auth.MiddlewareFunc())
	{
		_sessions.GET("", x.SessionController.Lists)
		_sessions.DELETE(":uid", x.SessionController.Remove)
		_sessions.DELETE("", x.SessionController.Clear)
	}

	_dsl := h.Group("dsl/:model", auth.MiddlewareFunc())
	{
		_dsl.POST("", x.DslController.Create)
		_dsl.POST("bulk-create", x.DslController.BulkCreate)
		_dsl.GET("_size", x.DslController.Size)
		_dsl.GET("", x.DslController.Find)
		_dsl.GET("_one", x.DslController.FindOne)
		_dsl.GET(":id", x.DslController.FindById)
		_dsl.PATCH("", x.DslController.Update)
		_dsl.PATCH(":id", x.DslController.UpdateById)
		_dsl.PUT(":id", x.DslController.Replace)
		_dsl.DELETE(":id", x.DslController.Delete)
		_dsl.POST("bulk-delete", x.DslController.BulkDelete)
		_dsl.POST("sort", x.DslController.Sort)
	}

	_pages := h.Group("pages", auth.MiddlewareFunc())
	{
		_pages.GET(":id", x.PagesController.GetOne)
	}

	return
}

// Auth 认证
func (x *API) Auth() (*jwt.HertzJWTMiddleware, error) {
	return jwt.New(&jwt.HertzJWTMiddleware{
		Realm:             x.Values.Namespace,
		Key:               []byte(x.Values.Key),
		Timeout:           time.Hour,
		SendAuthorization: false,
		SendCookie:        true,
		CookieMaxAge:      -1,
		SecureCookie:      true,
		CookieHTTPOnly:    true,
		CookieName:        "access_token",
		CookieSameSite:    http.SameSiteStrictMode,
		Authenticator: func(ctx context.Context, c *app.RequestContext) (_ interface{}, err error) {
			var dto struct {
				// 唯一标识，用户名或电子邮件
				Identity string `json:"identity,required" vd:"len($)>=4 || email($)"`
				// 密码
				Password string `json:"password,required" vd:"len($)>=8"`
			}
			if err = c.BindAndValidate(&dto); err != nil {
				c.Error(err)
				return
			}

			data, err := x.IndexService.Login(ctx, dto.Identity, dto.Password)
			if err != nil {
				c.Error(err)
				return
			}

			c.Set("identity", data)
			return data, nil
		},
		PayloadFunc: func(data interface{}) (claims jwt.MapClaims) {
			v := data.(common.Active)
			return jwt.MapClaims{
				"uid": v.UID,
				"jti": v.JTI,
			}
		},
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
			data := common.GetActive(c)
			if err := x.IndexService.LoginSession(ctx, data.UID, data.JTI); err != nil {
				c.Error(err)
				return
			}
			c.Status(http.StatusNoContent)
		},
		MaxRefresh: time.Hour,
		RefreshResponse: func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
			c.Status(http.StatusNoContent)
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.Error(errors.NewPublic(message).
				SetMeta(map[string]interface{}{
					"statusCode": http.StatusUnauthorized,
				}),
			)
		},
		TokenLookup: "cookie: access_token",
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			data := jwt.ExtractClaims(ctx, c)
			return common.Active{
				JTI: data["jti"].(string),
				UID: data["uid"].(string),
			}
		},
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			identity := data.(common.Active)
			if err := x.IndexService.AuthVerify(ctx, identity.UID, identity.JTI); err != nil {
				c.Error(err)
				return false
			}
			return true
		},
		LogoutResponse: func(ctx context.Context, c *app.RequestContext, code int) {
			data := common.GetActive(c)
			if err := x.IndexService.LogoutSession(ctx, data.UID); err != nil {
				c.Error(err)
				return
			}
			c.Status(http.StatusNoContent)
		},
	})
}

// AccessLog 日志
func (x *API) AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		end := time.Now()
		latency := end.Sub(start).Microseconds
		x.Transfer.Publish(context.Background(), "access_log", transfer.Payload{
			Tags: map[string]string{
				"method": string(c.Request.Header.Method()),
				"host":   string(c.Request.Host()),
				"path":   string(c.Request.Path()),
				"status": strconv.Itoa(c.Response.StatusCode()),
				"ip":     c.ClientIP(),
			},
			Fields: map[string]interface{}{
				"user_agent": string(c.Request.Header.UserAgent()),
				"query":      c.Request.QueryString(),
				"body":       string(c.Request.Body()),
				"cost":       latency(),
			},
			Time: start,
		})
	}
}

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
				"message": e.Description(),
			})
			break
		case *binding.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": e.Error(),
			})
			break
		case *validator.Error:
			c.JSON(http.StatusBadRequest, utils.H{
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
func (x *API) Initialize(ctx context.Context) (err error) {
	// 加载自定义验证
	validation.Extend()

	// 传输指标
	if err = x.Transfer.Set(transfer.Option{
		Measurement: "access_log",
	}); err != nil {
		return
	}

	// 订阅动态配置
	if err = x.ValuesService.Sync(ctx); err != nil {
		return
	}

	return
}
