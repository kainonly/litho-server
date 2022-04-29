package system

import (
	"api/app/departments"
	"api/app/pages"
	"api/app/roles"
	"api/app/users"
	"api/app/vars"
	"api/common"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thoas/go-funk"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	Service     *Service
	Users       *users.Service
	Roles       *roles.Service
	Departments *departments.Service
	Pages       *pages.Service
	Passport    *passport.Passport
	Vars        *vars.Service
}

func (x *Controller) Index(c *gin.Context) interface{} {
	return gin.H{
		"time": time.Now(),
	}
}

// AuthLogin 登录
func (x *Controller) AuthLogin(c *gin.Context) interface{} {
	var body struct {
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	// 用户验证
	data, err := x.Users.FindOneByUsernameOrEmail(ctx, body.User)
	if err != nil {
		c.Set("code", "AUTH_INCORRECT")
		return err
	}
	if err := helper.PasswordVerify(body.Password, data.Password); err != nil {
		c.Set("code", "AUTH_INCORRECT")
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
	if err := x.Service.SetSession(ctx, data.ID.Hex(), jti); err != nil {
		return err
	}
	// 写入日志
	dto := common.NewLoginLogV10(data, jti, c.ClientIP(), c.Request.UserAgent())
	go x.Service.WriteLoginLog(context.TODO(), dto)
	// 返回
	c.SetCookie("access_token", ts, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

// AuthVerify 主动验证
func (x *Controller) AuthVerify(c *gin.Context) interface{} {
	ts, err := c.Cookie("access_token")
	if err != nil {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	if _, err = x.Passport.Verify(ts); err != nil {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return err
	}
	return nil
}

// AuthCode 申请刷新验证码
func (x *Controller) AuthCode(c *gin.Context) interface{} {
	claims, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	jti := claims.(jwt.MapClaims)["jti"].(string)
	code := funk.RandomString(8)
	ctx := c.Request.Context()
	if err := x.Service.CreateCode(ctx, jti, code, time.Minute); err != nil {
		return err
	}
	return gin.H{"code": code}
}

// AuthRefresh 刷新认证
func (x *Controller) AuthRefresh(c *gin.Context) interface{} {
	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	// 获取载荷
	value, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	claims := value.(jwt.MapClaims)
	jti := claims["jti"].(string)
	ctx := c.Request.Context()
	// 刷新验证
	result, err := x.Service.VerifyCode(ctx, jti, body.Code)
	if err != nil {
		return err
	}
	if !result {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	if err = x.Service.DeleteCode(ctx, jti); err != nil {
		return err
	}
	// 继承 jti 创建新 Token
	var ts string
	if ts, err = x.Passport.Create(jti,
		claims["context"].(map[string]interface{}),
	); err != nil {
		return err
	}
	c.SetCookie("access_token", ts, 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

// AuthLogout 登出
func (x *Controller) AuthLogout(c *gin.Context) interface{} {
	c.SetCookie("access_token", "", 0, "", "", true, true)
	c.SetSameSite(http.SameSiteStrictMode)
	return nil
}

// CaptchaUser 用户验证码
func (x *Controller) CaptchaUser(c *gin.Context) interface{} {
	var query struct {
		Email string `form:"email" binding:"required,email"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Users.FindOneByEmail(ctx, query.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("该用户邮箱不存在")
		}
		return err
	}
	exists, err := x.Service.ExistsCode(ctx, query.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("您已获取验证码，请稍后再试~")
	}
	code := funk.RandomString(8)
	if err = x.Service.CreateCode(ctx, query.Email, code, time.Minute*5); err != nil {
		return err
	}
	if err = x.Service.EmailCode(data.Username, code, []string{query.Email}); err != nil {
		return err
	}
	return nil
}

// VerifyUser 校验验证码
func (x *Controller) VerifyUser(c *gin.Context) interface{} {
	var body struct {
		Email   string `json:"email" binding:"required,email"`
		Captcha string `json:"captcha" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.VerifyCode(ctx, body.Email, body.Captcha)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("您的验证码不正确")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 5).Unix(),
		"iss": body.Email,
	})
	ts, err := token.SignedString([]byte(x.Service.Values.Key))
	if err != nil {
		return err
	}
	return gin.H{
		"token": ts,
	}
}

// ResetUser 重置用户密码
func (x *Controller) ResetUser(c *gin.Context) interface{} {
	var body struct {
		Token    string `json:"token" binding:"required,jwt"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	token, err := jwt.Parse(body.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("验证失败，签名方式不一致")
		}
		return []byte(x.Service.Values.Key), nil
	})
	if err != nil {
		return err
	}
	ctx := c.Request.Context()
	password, _ := helper.PasswordHash(body.Password)
	email := token.Claims.(jwt.MapClaims)["iss"].(string)
	if err = x.Users.UpdateOneByEmail(ctx, email, bson.M{
		"$set": bson.M{
			"password": password,
		},
	}); err != nil {
		return err
	}
	return nil
}

// CheckUser 检查当前用户可变更属性
func (x *Controller) CheckUser(c *gin.Context) interface{} {
	claims, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	ctx := c.Request.Context()
	userId, _ := primitive.ObjectIDFromHex(claims.(jwt.MapClaims)["context"].(map[string]interface{})["uid"].(string))
	var query struct {
		Key   string `form:"key"`
		Value string `form:"value"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		return err
	}
	var count int64
	var err error
	switch query.Key {
	case "username":
		count, err = x.Users.Count(ctx, bson.M{
			"$and": bson.A{
				bson.M{"_id": bson.M{"$ne": userId}},
				bson.M{"username": query.Value},
			},
		})
		break
	case "email":
		count, err = x.Users.Count(ctx, bson.M{
			"$and": bson.A{
				bson.M{"_id": bson.M{"$ne": userId}},
				bson.M{"email": bson.M{"$ne": ""}},
				bson.M{"email": query.Value},
			},
		})
		break
	}
	if err != nil {
		return err
	}
	c.Header("wpx-exists", strconv.FormatBool(count != 0))
	return nil
}

// GetUser 获取用户信息
func (x *Controller) GetUser(c *gin.Context) interface{} {
	claims, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	ctx := c.Request.Context()
	userId, _ := primitive.ObjectIDFromHex(claims.(jwt.MapClaims)["context"].(map[string]interface{})["uid"].(string))
	var data common.User
	if err := x.Users.FindOneById(ctx, userId, &data); err != nil {
		return err
	}
	result := gin.H{
		"username":    data.Username,
		"email":       data.Email,
		"name":        data.Name,
		"avatar":      data.Avatar,
		"feishu":      data.Feishu,
		"sessions":    data.Sessions,
		"last":        data.Last,
		"create_time": data.CreateTime,
	}
	var err error
	if result["roles"], err = x.Roles.FindNamesById(ctx, data.Roles); err != nil {
		return err
	}
	if data.Department != nil {
		if result["department"], err = x.Departments.FindNameById(ctx, *data.Department); err != nil {
			return err
		}
	}
	return result
}

// SetUser 设置用户信息
func (x *Controller) SetUser(c *gin.Context) interface{} {
	var headers struct {
		Action string `header:"wpx-action"`
	}
	if err := c.ShouldBindHeader(&headers); err != nil {
		return err
	}
	claims, exists := c.Get(common.TokenClaimsKey)
	if !exists {
		c.Set("status_code", 401)
		c.Set("code", "AUTH_EXPIRED")
		return common.AuthExpired
	}
	ctx := c.Request.Context()
	userId, _ := primitive.ObjectIDFromHex(claims.(jwt.MapClaims)["context"].(map[string]interface{})["uid"].(string))
	switch headers.Action {
	case "profile":
		var body struct {
			Username string `json:"username,omitempty" bson:"username,omitempty"`
			Name     string `json:"name" binding:"required"`
			Avatar   string `json:"avatar" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return err
		}
		if err := x.Users.UpdateOneById(ctx, userId, bson.M{
			"$set": body,
		}); err != nil {
			return err
		}
		if body.Username != "" {
			if err := x.AuthLogout(c); err != nil {
				return err
			}
		}
		break
	case "password":
		var body struct {
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return err
		}
		if body.Password != "" {
			body.Password, _ = helper.PasswordHash(body.Password)
		}
		if err := x.Users.UpdateOneById(ctx, userId, bson.M{
			"$set": bson.M{
				"password": body.Password,
			},
		}); err != nil {
			return err
		}
		break
	case "email":
		var body struct {
			Email string `json:"email" binding:"omitempty,email"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return err
		}
		if err := x.Users.UpdateOneById(ctx, userId, bson.M{
			"$set": bson.M{
				"email": body.Email,
			},
		}); err != nil {
			return err
		}
	case "unlink":
		var body struct {
			Type string `json:"type" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return err
		}
		if err := x.Users.UpdateOneById(ctx, userId, bson.M{
			"$unset": bson.M{
				body.Type: "",
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

// GetSessions 获取会话
func (x *Controller) GetSessions(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	values, err := x.Service.GetSessions(ctx)
	if err != nil {
		return err
	}
	return values
}

// DeleteSessions 删除所有会话
func (x *Controller) DeleteSessions(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	if err := x.Service.DeleteSessions(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteSession 删除会话
func (x *Controller) DeleteSession(c *gin.Context) interface{} {
	var uri struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	if err := x.Service.DeleteSession(ctx, uri.Id); err != nil {
		return err
	}
	return nil
}

// Sort 通用排序
func (x *Controller) Sort(c *gin.Context) interface{} {
	var uri struct {
		Model string `uri:"model"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	var body struct {
		Sort []primitive.ObjectID `json:"sort" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	ctx := c.Request.Context()
	result, err := x.Service.Sort(ctx, uri.Model, body.Sort)
	if err != nil {
		return err
	}
	return result
}

// Navs 页面导航
func (x *Controller) Navs(c *gin.Context) interface{} {
	ctx := c.Request.Context()
	navs, err := x.Pages.Navs(ctx)
	if err != nil {
		return err
	}
	return navs
}

func (x *Controller) Dynamic(c *gin.Context) interface{} {
	var uri struct {
		Id string `uri:"id" binding:"required,objectId"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		return err
	}
	ctx := c.Request.Context()
	data, err := x.Pages.FindOneById(ctx, uri.Id)
	if err != nil {
		return err
	}
	return data
}
