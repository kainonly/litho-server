package feishu

import (
	"api/app/system"
	"api/app/users"
	"api/common"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"net/http"
)

type Controller struct {
	Service  *Service
	System   *system.Service
	Users    *users.Service
	Passport *passport.Passport
}

func (x *Controller) Challenge(c *gin.Context) interface{} {
	var body struct {
		Encrypt string `json:"encrypt"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	content, err := x.Service.Decrypt(body.Encrypt, x.Service.Values.Feishu.EncryptKey)
	if err != nil {
		return err
	}
	var dto struct {
		Challenge string `json:"challenge"`
		Token     string `json:"token"`
		Type      string `json:"type"`
	}
	if err = jsoniter.Unmarshal([]byte(content), &dto); err != nil {
		return err
	}
	if dto.Token != x.Service.Values.Feishu.VerificationToken {
		return errors.New("验证令牌不一致")
	}
	return gin.H{
		"challenge": dto.Challenge,
	}
}

func (x *Controller) T(c *gin.Context) interface{} {
	token, err := x.Service.GetTenantAccessToken(c.Request.Context())
	if err != nil {
		return err
	}
	return token
}

func (x *Controller) OAuth(c *gin.Context) interface{} {
	var query struct {
		Code  string `form:"code" binding:"required"`
		State string `form:"state"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	user, err := x.Service.GetAccessToken(ctx, query.Code)
	if err != nil {
		return err
	}
	data, err := x.Users.FindOneByFeishu(ctx, user.OpenId)
	if err != nil {
		return err
	}
	// 创建 Token
	jti := helper.Uuid()
	var ts string
	if ts, err = x.Passport.Create(jti, gin.H{
		"uid": data.ID.Hex(),
	}); err != nil {
		return err
	}
	// 设置会话
	if err := x.System.SetSession(ctx, data.ID.Hex(), jti); err != nil {
		return err
	}
	// 写入日志
	dto := common.NewLoginLogV10(data, jti, c.ClientIP(), c.Request.UserAgent())
	go x.System.WriteLoginLog(context.TODO(), dto)
	// 返回
	c.SetCookie("access_token", ts, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	c.Redirect(302, "https://xconsole.kainonly.com:8443/")
	return nil
}
