package feishu

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"server/app/system"
	"server/app/users"
	"server/common"
	"server/model"
)

type Controller struct {
	*common.Inject
	Feishu   *Service
	Users    *users.Service
	System   *system.Service
	Passport *passport.Passport
}

// Challenge 事件订阅
func (x *Controller) Challenge(c *gin.Context) interface{} {
	var body struct {
		Encrypt string `json:"encrypt"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	content, err := x.Feishu.Decrypt(body.Encrypt, x.Values.FeishuEncryptKey)
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
	if dto.Token != x.Values.FeishuVerificationToken {
		return errors.New("验证令牌不一致")
	}
	return gin.H{
		"challenge": dto.Challenge,
	}
}

type State struct {
	Action string `json:"action,omitempty"`
}

// OAuth 第三方登录与关联
func (x *Controller) OAuth(c *gin.Context) interface{} {
	var query struct {
		Code  string `form:"code" binding:"required"`
		State string `form:"state"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	userData, err := x.Feishu.GetAccessToken(ctx, query.Code)
	if err != nil {
		return err
	}
	var state State
	if err = jsoniter.Unmarshal([]byte(query.State), &state); err != nil {
		return err
	}
	switch state.Action {
	case "link":
		ts, err := c.Cookie("access_token")
		if err != nil {
			c.Set("status_code", 401)
			c.Set("code", "AUTH_EXPIRED")
			return common.AuthExpired
		}
		var claims jwt.MapClaims
		if claims, err = x.Passport.Verify(ts); err != nil {
			c.Set("status_code", 401)
			c.Set("code", "AUTH_EXPIRED")
			return err
		}
		userId, _ := primitive.ObjectIDFromHex(claims["context"].(map[string]interface{})["uid"].(string))
		if err := x.Users.UpdateOneById(ctx, userId, bson.M{
			"$set": bson.M{
				"feishu": userData,
			},
		}); err != nil {
			return err
		}
		c.Redirect(302, fmt.Sprintf(`%s/#/authorized`, x.Values.Console))
		return nil
	}
	var user model.User
	if err := x.Users.FindOneByFeishu(ctx, userData.OpenId, &user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Redirect(302, fmt.Sprintf(`%s/#/unauthorize`, x.Values.Console))
		}
		return err
	}

	// 创建 Token
	jti := helper.Uuid()
	var ts string
	if ts, err = x.Passport.Create(jti, gin.H{
		"uid": user.ID.Hex(),
	}); err != nil {
		return err
	}
	// 设置会话
	if err := x.System.SetSession(ctx, user.ID.Hex(), jti); err != nil {
		return err
	}
	// 写入日志
	ip := c.GetHeader("X-Forwarded-For")
	logData := model.NewLoginLogV10(user, jti, ip, c.Request.UserAgent())
	go x.System.PushLoginLog(context.TODO(), logData)
	// 返回
	c.SetCookie("access_token", ts, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	c.Redirect(302, fmt.Sprintf(`%s/`, x.Values.Console))
	return nil
}
