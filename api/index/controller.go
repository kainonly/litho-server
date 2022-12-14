package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/weplanx/server/common"
	"github.com/weplanx/utils/passlib"
	"net/http"
	"reflect"
	"strings"
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
		// 电子邮件
		Email string `json:"email,required" vd:"email($)"`
		// 密码
		Password string `json:"password,required" vd:"len($)>=8"`
	}
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	ts, err := x.IndexService.Login(ctx, dto.Email, dto.Password)
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
	claims := common.GetClaims(c)
	code, err := x.IndexService.GetRefreshCode(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code": code,
	})
}

type RefreshTokenDto struct {
	Code string `json:"code,required"`
}

// RefreshToken 刷新令牌
// @router /refresh_token [POST]
func (x *Controller) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var dto RefreshTokenDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.GetClaims(c)
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
	claims := common.GetClaims(c)
	if err := x.IndexService.Logout(ctx, claims.UserId); err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "认证已注销",
	})
}

// GetOptions 返回通用配置
// @router /options [GET]
//func (x *Controller) GetOptions(ctx context.Context, c *app.RequestContext) {
//	var dto struct {
//		// 类型
//		Type string `query:"type,required"`
//	}
//	if err := c.BindAndValidate(&dto); err != nil {
//		c.Error(err)
//		return
//	}
//
//	data := x.IndexService.GetOptions(dto.Type)
//	if data == nil {
//		c.Error(errors.NewPublic("配置类型不存在"))
//		return
//	}
//
//	c.JSON(http.StatusOK, data)
//}

// GetUser 获取授权用户信息
// @router /user [GET]
func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	claims := common.GetClaims(c)
	data, err := x.IndexService.GetUser(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, data)
}

type SetUserDto struct {
	// 更新字段
	Set string `json:"$set,requred" vd:"in($, 'Email', 'Name', 'Avatar', 'Password')"`
	// 电子邮件
	Email string `json:"email,omitempty" vd:"(Set)$!='Email' || email($);msg:'必须是电子邮件'"`
	// 称呼
	Name string `json:"name,omitempty"`
	// 头像
	Avatar string `json:"avatar,omitempty"`
	// 密码
	Password string `json:"password,omitempty" vd:"(Set)$!='Password' || len($)>8;msg:'密码必须大于8位'"`
}

// SetUser 设置授权用户信息
// @router /user [POST]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto SetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	data := make(map[string]interface{})
	key := strings.ToLower(dto.Set)
	value := reflect.ValueOf(dto).FieldByName(dto.Set).Interface()
	if key == "password" {
		data[key], _ = passlib.Hash(value.(string))
	} else {
		data[key] = value
	}

	claims := common.GetClaims(c)
	if err := x.IndexService.SetUser(ctx, claims.UserId, data); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "设置成功",
	})
}
