package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/weplanx/server/utils/passlib"
	"github.com/weplanx/server/utils/passport"
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
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

// Login 登录认证
// @router /login [POST]
func (x *Controller) Login(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 唯一标识，用户名或电子邮件
		Identity string `json:"identity,required" vd:"len($)>=4 || email($)"`
		// 密码
		Password string `json:"password,required" vd:"len($)>=8"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	ts, err := x.IndexService.Login(ctx, dto.Identity, dto.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", ts, 0, "/", "", protocol.CookieSameSiteStrictMode, true, true)
	c.JSON(200, utils.H{
		"code":    0,
		"message": "登录认证成功",
	})
}

// Verify 主动验证
// @router /verify [GET]
func (x *Controller) Verify(ctx context.Context, c *app.RequestContext) {
	ts := c.Cookie("access_token")
	if ts == nil {
		c.JSON(401, utils.H{
			"code":    0,
			"message": "认证已失效请重新登录",
		})
		return
	}

	if _, err := x.IndexService.Verify(ctx, string(ts)); err != nil {
		c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
		c.JSON(401, utils.H{
			"code":    0,
			"message": "认证已失效请重新登录",
		})
		return
	}

	c.JSON(200, utils.H{
		"code":    0,
		"message": "验证成功",
	})
}

// GetRefreshCode 获取刷新令牌验证码
// @router /code [GET]
func (x *Controller) GetRefreshCode(ctx context.Context, c *app.RequestContext) {
	claims := passport.GetClaims(c)
	code, err := x.IndexService.GetRefreshCode(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": code,
	})
}

// RefreshToken 刷新令牌
// @router /refresh_token [POST]
func (x *Controller) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		Code string `json:"code,required"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := passport.GetClaims(c)
	ts, err := x.IndexService.RefreshToken(ctx, claims, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", ts, 0, "/", "", protocol.CookieSameSiteStrictMode, true, true)
	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "令牌刷新成功",
	})
}

// Logout 注销认证
// @router /logout [POST]
func (x *Controller) Logout(ctx context.Context, c *app.RequestContext) {
	claims := passport.GetClaims(c)
	if err := x.IndexService.Logout(ctx, claims.UserId); err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
	c.JSON(200, utils.H{
		"code":    0,
		"message": "认证已注销",
	})
}

// GetNavs 导航数据
// @router /navs [GET]
func (x *Controller) GetNavs(ctx context.Context, c *app.RequestContext) {
	claims := passport.GetClaims(c)
	data, err := x.IndexService.GetNavs(ctx, claims.UserId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetOptions 返回通用配置
// @router /options [GET]
func (x *Controller) GetOptions(ctx context.Context, c *app.RequestContext) {
	var dto struct {
		// 类型
		Type string `query:"type,required"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data := x.IndexService.GetOptions(dto.Type)
	if data == nil {
		c.Error(errors.NewPublic("配置类型不存在"))
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetUser 获取授权用户信息
// @router /user [GET]
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	claims := passport.GetClaims(c)
	data, err := x.IndexService.GetUser(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

type SetUserDto struct {
	// 用户名
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	// 电子邮件
	Email string `json:"email,omitempty" bson:"email,omitempty" vd:"$=='' || email($)"`
	// 称呼
	Name string `json:"name" bson:"name,omitempty"`
	// 头像
	Avatar string `json:"avatar" bson:"avatar,omitempty"`
	// 密码
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	// 重置
	Reset string `json:"reset,omitempty" vd:"in($, 'feishu')" bson:"reset"`
	// 更新时间
	UpdateTime time.Time `json:"-" bson:"update_time"`
}

// SetUser 设置授权用户信息
// @router /user [PATCH]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto SetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	// 密码转散列
	if dto.Password != "" {
		dto.Password, _ = passlib.Hash(dto.Password)
	}

	dto.UpdateTime = time.Now()
	claims := passport.GetClaims(c)
	if _, err := x.IndexService.SetUser(ctx, claims.UserId, dto); err != nil {
		c.Error(err)
		return
	}

	// 用户名变更，注销登录状态
	if dto.Username != "" {
		c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
		if err := x.IndexService.Logout(ctx, claims.UserId); err != nil {
			c.Error(err)
			return
		}
	}

	c.Status(http.StatusNoContent)
}
