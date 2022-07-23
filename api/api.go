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
	"github.com/weplanx/server/api/values"
	"github.com/weplanx/server/common"
	"net/http"
	"os"
	"time"
)

type API struct {
	Values        *common.Values
	ValuesService *values.Service
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

func (x *API) ErrHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		err := c.Errors.Last()
		if err == nil {
			return
		}

		if err.IsType(errors.ErrorTypePublic) {
			result := utils.H{"message": err.Error()}
			if meta, ok := err.Meta.(map[string]interface{}); ok {
				if meta["code"] != nil {
					result["code"] = meta["code"]
				}
			}
			c.JSON(http.StatusBadRequest, result)
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
