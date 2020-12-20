package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/hash"
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
	admin := c.Cache.Admin.Get(body.Username)
	if admin == nil {
		return errors.New("user does not exist or has been frozen")
	}
	if ok, _ := hash.Verify(body.Password, admin["password"].(string)); !ok {
		return errors.New("user login password is incorrect")
	}
	if err := authx.Create(ctx, typ.Cookie{
		Name:     "system",
		MaxAge:   0,
		Path:     "",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, jwt.MapClaims{
		"username": admin["username"],
		"role":     admin["role"],
	}, c.Cache.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	if err := authx.Verify(ctx, typ.Cookie{Name: "system",
		MaxAge:   0,
		Path:     "",
		Domain:   "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}, c.Cache.RefreshToken); err != nil {
		return err
	}
	return true
}

func (c *Controller) Resource(ctx *gin.Context) interface{} {
	resource := c.Cache.Resource.Get()
	//roleKeyids := []string{"*"}
	//role, err := c.Cache.RoleGet(roleKeyids, "resource")
	//if err != nil {
	//	return res.Error(err)
	//}
	return resource
}
