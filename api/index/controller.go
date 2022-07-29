package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/common/passlib"
	"net/http"
	"time"
)

type Controller struct {
	IndexService *Service
}

// Index 入口
// @router / [GET]
func (x *Controller) Index(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"msg":  "hi",
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

// GetNavs 导航数据
// @router /navs [GET]
func (x *Controller) GetNavs(ctx context.Context, c *app.RequestContext) {
	active := common.GetActive(c)

	data, err := x.IndexService.GetNavs(ctx, active.UID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, data)
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
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	active := common.GetActive(c)
	data, err := x.IndexService.GetUser(ctx, active.UID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

// SetUser 设置授权用户信息
// @router /user [PATCH]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 用户名
		Username string `json:"username,omitempty" bson:"username,omitempty"`
		// 电子邮件
		Email string `json:"email,omitempty" vd:"$='' || email($)"`
		// 称呼
		Name string `json:"name" bson:"name,omitempty"`
		// 头像
		Avatar string `json:"avatar" bson:"avatar,omitempty"`
		// 密码
		Password string `json:"password,omitempty" bson:"avatar,omitempty"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 密码不为空转散列
	if dto.Password != "" {
		dto.Password, _ = passlib.Hash(dto.Password)
	}

	active := common.GetActive(c)
	if _, err := x.IndexService.SetUser(ctx, active.UID, dto); err != nil {
		c.Error(err)
		return
	}

	// 用户名变更，注销登录状态
	if dto.Username != "" {
		cookie := &protocol.Cookie{}
		cookie.SetKey("access_token")
		cookie.SetValue("")
		c.Response.Header.SetCookie(cookie)

		if err := x.IndexService.LogoutSession(ctx, active.UID); err != nil {
			c.Error(err)
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// UnsetUser 取消授权用户信息
// @router /unset-user [POST]
func (x *Controller) UnsetUser(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 用户名
		Type string `json:"type,required"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	active := common.GetActive(c)
	if _, err := x.IndexService.UnsetUser(ctx, active.UID, dto.Type); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
