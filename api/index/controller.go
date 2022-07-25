package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

type Controller struct {
	AppService *Service
}

// Index 入口
// @router / [GET]
func (x *Controller) Index(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"time": x.AppService.Index(),
		"ip":   c.ClientIP(),
	})
}

// GetRefreshCode 获取刷新随机验证
// @router /code [GET]
func (x *Controller) GetRefreshCode(ctx context.Context, c *app.RequestContext) {}

// GetUser 获取授权用户信息
// @router /user [GET]
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {}

// SetUser 设置授权用户信息
// @router /user [PATCH]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {}
