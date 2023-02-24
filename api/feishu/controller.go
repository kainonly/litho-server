package feishu

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/utils/passport"
	"net/http"
)

type Controller struct {
	IndexService  *index.Service
	FeishuService *Service
	Values        *common.Values
	Passport      *passport.Passport
}

type ChallengeDto struct {
	Encrypt string `json:"encrypt,required"`
}

// Challenge 事件订阅
func (x *Controller) Challenge(ctx context.Context, c *app.RequestContext) {
	var dto ChallengeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	raw, err := x.FeishuService.Decrypt(dto.Encrypt, x.Values.FeishuEncryptKey)
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
	if data.Token != x.Values.FeishuVerificationToken {
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

// OAuth 第三方登录与关联
func (x *Controller) OAuth(ctx context.Context, c *app.RequestContext) {
	var dto OAuthDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	userData, err := x.FeishuService.GetUserAccessToken(ctx, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	switch dto.State.Action {
	case "link":
		ts := c.Cookie("access_token")
		if ts == nil {
			c.JSON(401, utils.H{
				"code":    0,
				"message": index.MsgAuthenticationExpired,
			})
			return
		}
		claims, err := x.IndexService.Verify(ctx, string(ts))
		if err != nil {
			c.SetCookie("access_token", "", -1, "", "", protocol.CookieSameSiteLaxMode, true, true)
			c.JSON(401, utils.H{
				"code":    0,
				"message": index.MsgAuthenticationExpired,
			})
			return
		}

		if _, err = x.FeishuService.Link(ctx, claims.UserId, userData); err != nil {
			c.Error(err)
			return
		}

		c.Redirect(302, []byte(fmt.Sprintf(`%s/%s/#/authorized`, x.Values.BaseUrl, dto.State.Locale)))
		return
	}

	var metadata model.LoginMetadata
	metadata.Channel = "feishu"
	ts, err := x.FeishuService.Login(ctx, userData.OpenId, &metadata)
	if err != nil {
		c.Redirect(302, []byte(fmt.Sprintf(`%s/%s/#/unauthorize`, x.Values.BaseUrl, dto.State.Locale)))
		return
	}

	//metadata.Ip = c.ClientIP()
	//var data model.LoginData
	//data.UserAgent = string(c.UserAgent())
	//go func() {
	//	if err := x.IndexService.WriteLoginLog(ctx, metadata, data); err != nil {
	//		logger.Error(err)
	//		return
	//	}
	//}()

	c.SetCookie("access_token", ts, 0, "", "", protocol.CookieSameSiteLaxMode, true, true)
	c.Redirect(302, []byte(fmt.Sprintf(`%s/%s/#/`, x.Values.BaseUrl, dto.State.Locale)))
}

func (x *Controller) CreateTasks(ctx context.Context, c *app.RequestContext) {
	result, err := x.FeishuService.CreateTask(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (x *Controller) GetTasks(ctx context.Context, c *app.RequestContext) {
	result, err := x.FeishuService.GetTasks(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}
