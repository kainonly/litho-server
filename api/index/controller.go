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

// Ping
// @router / [GET]
func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

// User Login
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
		"message": "ok",
	})
}

// User Verify
// @router /verify [GET]
func (x *Controller) Verify(ctx context.Context, c *app.RequestContext) {
	ts := c.Cookie("access_token")
	if ts == nil {
		c.JSON(401, utils.H{
			"code":    0,
			"message": MsgAuthenticationExpired,
		})
		return
	}

	if _, err := x.IndexService.Verify(ctx, string(ts)); err != nil {
		c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
		c.JSON(401, utils.H{
			"code":    0,
			"message": MsgAuthenticationExpired,
		})
		return
	}

	c.JSON(200, utils.H{
		"code":    0,
		"message": "ok",
	})
}

// Get Token Captcha
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

// Refresh Token
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
		"message": "ok",
	})
}

// Logout
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
		"message": "ok",
	})
}

// Get User Info
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
	Set      string `json:"$set,requred" vd:"in($, 'email', 'name', 'avatar', 'password')"`
	Email    string `json:"email,omitempty" vd:"(Set)$!='Email' || email($);msg:'must be email'"`
	Name     string `json:"name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Password string `json:"password,omitempty" vd:"(Set)$!='Password' || len($)>8;msg:'must be greater than 8 characters'"`
}

// Set User Info
// @router /user [POST]
func (x *Controller) SetUser(ctx context.Context, c *app.RequestContext) {
	var dto SetUserDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	data := make(map[string]interface{})
	path := strings.ToUpper(dto.Set[:1]) + dto.Set[1:]
	value := reflect.ValueOf(dto).FieldByName(path).Interface()
	if dto.Set == "password" {
		data[dto.Set], _ = passlib.Hash(value.(string))
	} else {
		data[dto.Set] = value
	}

	claims := common.GetClaims(c)
	_, err := x.IndexService.SetUser(ctx, claims.UserId, data)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "ok",
	})
}
