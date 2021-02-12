package controller

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/hash"
	"github.com/kainonly/gin-extra/mvcx"
	"github.com/kainonly/gin-extra/typ"
	"lab-api/application/common"
	"net/http"
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
	if err := authx.Create(ctx, typ.Cookie{
		Name:     "system",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, jwt.MapClaims{
		"username": admin["username"],
		"role":     admin["role"],
	}, c.Redis.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	if err := authx.Verify(ctx, typ.Cookie{
		Name:     "system",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, c.Redis.RefreshToken); err != nil {
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
	resource := c.Redis.Resource.Get()
	return resource
}
