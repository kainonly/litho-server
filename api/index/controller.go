package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/server/common"
	"net/http"
	"time"
)

type Controller struct {
	IndexService *Service
	Csrf         *csrf.Csrf
}

func (x *Controller) Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, utils.H{
		"ip":   c.ClientIP(),
		"time": time.Now(),
	})
}

type LoginDto struct {
	Email    string `json:"email,required" vd:"email($)"`
	Password string `json:"password,required" vd:"len($)>=8"`
}

func (x *Controller) Login(ctx context.Context, c *app.RequestContext) {
	var dto LoginDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	ts, err := x.IndexService.Login(ctx, dto.Email, dto.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", ts, -1, "/", "", protocol.CookieSameSiteNoneMode, false, false)
	c.JSON(200, utils.H{
		"code":    0,
		"message": "ok",
	})
}

func (x *Controller) Verify(ctx context.Context, c *app.RequestContext) {
	ts := c.Cookie("access_token")
	if ts == nil {
		c.JSON(401, utils.H{
			"code":    0,
			"message": common.MsgAuthenticationExpired,
		})
		return
	}

	if _, err := x.IndexService.Verify(ctx, string(ts)); err != nil {
		c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteNoneMode, false, false)
		c.JSON(401, utils.H{
			"code":    0,
			"message": common.MsgAuthenticationExpired,
		})
		return
	}

	c.JSON(200, utils.H{
		"code":    0,
		"message": "ok",
	})
}

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

	c.SetCookie("access_token", ts, 0, "/", "", protocol.CookieSameSiteNoneMode, false, false)
	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "ok",
	})
}

func (x *Controller) Logout(ctx context.Context, c *app.RequestContext) {
	claims := common.GetClaims(c)
	if err := x.IndexService.Logout(ctx, claims.UserId); err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", protocol.CookieSameSiteNoneMode, false, false)
	c.JSON(http.StatusOK, utils.H{
		"code":    0,
		"message": "ok",
	})
}
