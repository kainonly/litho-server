package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/utils/errors"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

type API struct {
	Values *common.Values
}

// Engine 创建服务
func (x *API) Engine() *gin.Engine {
	r := gin.New()

	// 全局中间件
	if os.Getenv("GIN_MODE") == "release" {
		// 生产环境
		r.SetTrustedProxies(strings.Split(os.Getenv("TRUSTED_PROXIES"), ","))
		logger, _ := zap.NewProduction()
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	} else {
		// 开发环境
		r.SetTrustedProxies([]string{})
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     x.Values.Cors.AllowOrigins,
		AllowMethods:     x.Values.Cors.AllowMethods,
		AllowHeaders:     x.Values.Cors.AllowHeaders,
		ExposeHeaders:    x.Values.Cors.ExposeHeaders,
		AllowCredentials: x.Values.Cors.AllowCredentials,
		MaxAge:           time.Duration(x.Values.Cors.MaxAge) * time.Second,
	}))
	r.Use(errors.Handler())

	return r
}
