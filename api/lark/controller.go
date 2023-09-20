package lark

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"net/http"
	"time"
)

type Controller struct {
	V        *common.Values
	Passport *passport.Passport

	LarkService  *Service
	IndexService *index.Service
}

type ChallengeDto struct {
	Encrypt string `json:"encrypt,required"`
}

func (x *Controller) Challenge(ctx context.Context, c *app.RequestContext) {
	var dto ChallengeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	raw, err := x.LarkService.Decrypt(dto.Encrypt, x.V.LarkEncryptKey)
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
		c.Error(errors.NewPublic("The local configuration token does not match the authentication token"))
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"challenge": data.Challenge,
	})
}

type OAuthDto struct {
	Code  string   `query:"code,required"`
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
	userData, err := x.LarkService.GetUserAccessToken(ctx, dto.Code)
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
		if claims, err = x.IndexService.Verify(ctx, string(ts)); err != nil {
			common.ClearAccessToken(c)
			c.JSON(401, utils.H{
				"code":    0,
				"message": common.ErrAuthenticationExpired.Error(),
			})
			return
		}

		if _, err = x.LarkService.Link(ctx, claims.UserId, userData); err != nil {
			c.Error(err)
			return
		}

		c.Redirect(302, []byte(fmt.Sprintf(`%s/%s/#/authorized`, x.V.Console, dto.State.Locale)))
		return
	}

	var r *LoginResult
	if r, err = x.LarkService.Login(ctx, userData.OpenId); err != nil {
		c.Redirect(302, []byte(fmt.Sprintf(`%s/%s/#/unauthorize`, x.V.Console, dto.State.Locale)))
		return
	}

	go func() {
		data := model.NewLogsetLogined(
			r.User.ID, string(c.GetHeader(x.V.Ip)), "lark", string(c.UserAgent()))
		wctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err = x.IndexService.WriteLogsetLogined(wctx, data); err != nil {
			hlog.Fatal(err)
		}
	}()

	common.SetAccessToken(c, r.AccessToken)
	c.Redirect(302, []byte(fmt.Sprintf(`%s/%s`, x.V.Console, dto.State.Locale)))
}

func (x *Controller) CreateTasks(ctx context.Context, c *app.RequestContext) {
	r, err := x.LarkService.CreateTask(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (x *Controller) GetTasks(ctx context.Context, c *app.RequestContext) {
	r, err := x.LarkService.GetTasks(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}
