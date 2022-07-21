package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/cors"
	"github.com/weplanx/server/api/values"
	"github.com/weplanx/server/common"
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
	//r.Use(errs.Handler())

	// 订阅动态配置
	if err = x.ValuesService.Sync(context.TODO()); err != nil {
		return
	}

	return
}
