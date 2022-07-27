package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/server/api/pages"
	"github.com/weplanx/server/common"
	"net/http"
	"time"
)

type Controller struct {
	IndexService *Service
	PagesService *pages.Service
}

// Index 获取导航
// @router / [GET]
func (x *Controller) Index(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{})
}

// GetRefreshCode 获取刷新令牌验证码
// @router /code [GET]
func (x *Controller) GetRefreshCode(ctx context.Context, c *app.RequestContext) {
	active := common.GetActive(c)
	code, _ := gonanoid.Nanoid()
	if err := x.IndexService.CreateCaptcha(ctx, active.UID, code, 15*time.Second); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": code,
	})
}

// VerifyRefreshCode 校验刷新令牌验证码
// @router /refresh_token [POST]
func (x *Controller) VerifyRefreshCode(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Code string `json:"code,required"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	active := common.GetActive(c)
	if err := x.IndexService.VerifyCaptcha(ctx, active.UID, dto.Code); err != nil {
		c.Error(err)
		return
	}

	c.Next(ctx)
}

// GetUser 获取授权用户信息
// @router /user [GET]
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {}

// SetUser 设置授权用户信息
// @router /user [PATCH]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {}
