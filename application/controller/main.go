package controller

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/hash"
	"github.com/kainonly/gin-extra/mvcx"
	"lab-api/application/common"
	"lab-api/application/model"
)

type Controller struct {
	common.Dependency
}

type LoginBody struct {
	Username string `binding:"required,min=4,max=20"`
	Password string `binding:"required,min=12,max=20"`
}

func (c *Controller) Login(ctx *gin.Context) interface{} {
	var body LoginBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	admin := c.Redis.Admin.Get(body.Username)
	if admin == nil {
		return errors.New("user does not exist or has been frozen")
	}
	cb := context.Background()
	if !c.Redis.UserLock.Check(cb, "admin:"+body.Username) {
		c.Redis.UserLock.Lock(cb, "admin:"+body.Username)
		return mvcx.Response{
			Code: 2,
			Msg:  "you have failed to log in too many times, please try again later",
		}
	}
	if ok, _ := hash.Verify(body.Password, admin["password"].(string)); !ok {
		c.Redis.UserLock.Inc(cb, "admin:"+body.Username)
		return errors.New("user login password is incorrect")
	}
	c.Redis.UserLock.Remove("admin:" + body.Username)
	if err := authx.Create(ctx, common.SystemCookie, jwt.MapClaims{
		"user": admin["username"],
	}, c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	if err := authx.Verify(ctx, common.SystemCookie, c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Logout(ctx *gin.Context) interface{} {
	if err := authx.Destory(ctx, "system", c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Resource(ctx *gin.Context) interface{} {
	var err error
	var auth jwt.MapClaims
	if auth, err = authx.Get(ctx); err != nil {
		return err
	}
	userData := c.Redis.Admin.Get(auth["user"].(string))
	roleKeys := userData["role"].([]interface{})
	keys := make([]string, len(roleKeys))
	for index, value := range roleKeys {
		keys[index] = value.(string)
	}
	roleResource := c.Redis.Role.Get(keys, "resource")
	if userData["resource"] != nil {
		roleResource.Add(userData["resource"].([]interface{})...)
	}
	resource := c.Redis.Resource.Get()
	lists := make([]interface{}, 0)
	for _, item := range resource {
		if roleResource.Contains(item["key"].(string)) {
			lists = append(lists, item)
		}
	}
	return lists
}

func (c *Controller) Information(ctx *gin.Context) interface{} {
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	data := make(map[string]interface{})
	c.Db.Model(&model.Admin{}).
		Where("username = ?", auth.(jwt.MapClaims)["user"]).
		First(&data)
	return data
}

type updateBody struct {
	OldPassword string
	NewPassword string
	Email       string
	Phone       string
	Call        string
	Avatar      string
}

func (c *Controller) Update(ctx *gin.Context) interface{} {
	var body updateBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	username := auth.(jwt.MapClaims)["user"]
	data := model.Admin{
		Email:  body.Email,
		Phone:  body.Phone,
		Call:   body.Call,
		Avatar: body.Avatar,
	}
	if body.NewPassword != "" || body.OldPassword != "" {
		if body.NewPassword != body.OldPassword {
			return errors.New("the old and new password verification is inconsistent")
		}
		var adminData model.Admin
		c.Db.Model(&model.Admin{}).
			Where("username = ?", username).
			First(&adminData)
		if ok, _ := hash.Verify(body.OldPassword, adminData.Password); !ok {
			return mvcx.Response{
				Code: 2,
				Msg:  "password verification failed",
			}
		}
		data.Password, _ = hash.Make(body.NewPassword)
	}
	c.Db.Model(&model.Admin{}).
		Where("username = ?", username).
		Updates(data)
	c.Redis.Admin.Clear()
	return true
}
