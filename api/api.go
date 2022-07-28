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
	"github.com/weplanx/server/common/validation"
	"net/http"
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
	*common.Inject

	Hertz *server.Hertz

	IndexController   *index.Controller
	IndexService      *index.Service
	ValuesController  *values.Controller
	ValuesService     *values.Service
	SessionController *sessions.Controller
	SessionService    *sessions.Service
	DslController     *dsl.Controller
	DslService        *dsl.Service
	UsersController   *users.Controller
	UsersService      *users.Service
}

func (x *API) Routes() (h *server.Hertz, err error) {
	h = x.Hertz
	h.Use(x.ErrHandler())

	var auth *jwt.HertzJWTMiddleware
	if auth, err = x.Auth(); err != nil {
		return
	}

	h.POST("login", auth.LoginHandler)

	api := h.Group("", auth.MiddlewareFunc())
	{
		api.GET("", x.IndexController.Index)
		api.GET("code", x.IndexController.GetRefreshCode)
		api.POST("refresh_token", x.IndexController.VerifyRefreshCode, auth.RefreshHandler)

		api.GET("user", x.IndexController.GetUser)
		api.PATCH("user", x.IndexController.SetUser)
		api.DELETE("user", auth.LogoutHandler)

		x.ValuesController.In(api.Group("values"))
		x.SessionController.In(api.Group("sessions"))
		x.DslController.In(api.Group("dsl/:model"))
	}

	return
}

// Auth 认证
func (x *API) Auth() (*jwt.HertzJWTMiddleware, error) {
	return jwt.New(&jwt.HertzJWTMiddleware{
		Realm:             x.Values.App.Namespace,
		Key:               []byte(x.Values.App.Key),
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
				Identity string `json:"identity,required" vd:"len($)>=4 || email($)"`
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

		switch any := err.Err.(type) {
		case decoder.SyntaxError:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": any.Description(),
			})
			break
		case *binding.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": any.Error(),
			})
			break
		case *validator.Error:
			c.JSON(http.StatusBadRequest, utils.H{
				"message": any.Error(),
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

	// 订阅动态配置
	if err = x.ValuesService.Sync(ctx); err != nil {
		return
	}

	return
}
