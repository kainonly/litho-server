package lark

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
)

type Controller struct {
	V        *common.Values
	Passport *passport.Passport

	LarkX  *Service
	IndexX *index.Service
}

type ChallengeDto struct {
	Encrypt string `json:"encrypt" vd:"required"`
}

func (x *Controller) Challenge(ctx context.Context, c *app.RequestContext) {
	var dto ChallengeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	raw, err := x.LarkX.Decrypt(dto.Encrypt, x.V.LarkEncryptKey)
	if err != nil {
		c.Error(err)
		return
	}
	var data struct {
		Challenge string `json:"challenge"`
		Token     string `json:"token"`
		Type      string `json:"type"`
	}
	if err = sonic.UnmarshalString(raw, &data); err != nil {
		c.Error(err)
		return
	}
	if data.Token != x.V.LarkVerificationToken {
		c.Error(help.E(
			"lark.VerificationTokenNotMatch",
			"the local configuration token does not match the authentication token"),
		)
		return
	}

	c.JSON(200, utils.H{
		"challenge": data.Challenge,
	})
}

type OAuthDto struct {
	Code  string   `query:"code" vd:"required"`
	State StateDto `query:"state"`
}

type StateDto struct {
	Action string `json:"action,omitempty"`
	Locale string `json:"locale,omitempty"`
}

func (x *Controller) OAuth(ctx context.Context, c *app.RequestContext) {
	var dto OAuthDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	userData, err := x.LarkX.GetUserAccessToken(ctx, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	switch dto.State.Action {
	case "link":
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.JSON(401, utils.H{
				"code":    0,
				"message": common.ErrAuthenticationExpired.Error(),
			})
			return
		}
		var claims passport.Claims
		if claims, err = x.IndexX.Verify(ctx, string(ts)); err != nil {
			common.ClearAccessToken(c)
			c.JSON(401, help.E(
				"lark.AuthenticationExpired",
				common.ErrAuthenticationExpired.Error(),
			))
			return
		}

		if _, err = x.LarkX.Link(ctx, claims.UserId, userData); err != nil {
			c.Error(err)
			return
		}
		c.Redirect(302, x.RedirectUrl("/#/authorized", dto.State.Locale))
		return
	}

	var r *LoginResult
	if r, err = x.LarkX.Login(ctx, userData.OpenId); err != nil {
		c.Redirect(302, x.RedirectUrl("/#/unauthorize", dto.State.Locale))
		return
	}

	data := model.NewLogsetLogin(r.User.ID, string(c.GetHeader(x.V.Ip)), "lark", string(c.UserAgent()))
	if err = x.IndexX.WriteLogsetLogin(ctx, data); err != nil {
		c.Error(err)
		return
	}

	common.SetAccessToken(c, r.AccessToken)
	c.Redirect(302, x.RedirectUrl("", dto.State.Locale))
}

func (x *Controller) RedirectUrl(path string, locale string) []byte {
	u := x.V.Console
	if locale != "" {
		u += "/" + locale
	}
	return []byte(u + path)
}

func (x *Controller) CreateTasks(ctx context.Context, c *app.RequestContext) {
	r, err := x.LarkX.CreateTask(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, r)
}

func (x *Controller) GetTasks(ctx context.Context, c *app.RequestContext) {
	r, err := x.LarkX.GetTasks(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}
