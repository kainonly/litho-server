package api

import (
	"context"
	"github.com/bytedance/go-tagexpr/v2/binding"
	"github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/bytedance/sonic/decoder"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/jwt"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/api/users"
	"github.com/weplanx/server/api/values"
	"github.com/weplanx/server/common"
	"net/http"
	"os"
	"time"
)

type API struct {
	Values        *common.Values
	ValuesService *values.Service
	IndexService  *index.Service
	UsersService  *users.Service
}

// Engine 创建服务
func (x *API) Engine() (h *server.Hertz, err error) {
	opts := []config.Option{
		server.WithHostPorts(":3000"),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	h = server.Default(opts...)
	// 全局中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins:     x.Values.Cors.AllowOrigins,
		AllowMethods:     x.Values.Cors.AllowMethods,
		AllowHeaders:     x.Values.Cors.AllowHeaders,
		AllowCredentials: x.Values.Cors.AllowCredentials,
		ExposeHeaders:    x.Values.Cors.ExposeHeaders,
		MaxAge:           time.Duration(x.Values.Cors.MaxAge) * time.Second,
	}))
	h.Use(x.ErrHandler())

	// 订阅动态配置
	if err = x.ValuesService.Sync(context.TODO()); err != nil {
		return
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
