package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/hash"
	"github.com/kainonly/gin-extra/tokenx"
	"github.com/kainonly/gin-extra/typ"
	"lab-api/application/common"
	"net/http"
	"time"
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
	authx.Create(ctx, typ.Cookie{
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
	}, c.Cache.RefreshToken)
	return true
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	tokenString, _ := ctx.Cookie("system_token")
	tokenx.Verify(tokenString, func(claims jwt.MapClaims) (jwt.MapClaims, error) {
		var err error
		jti := claims["jti"].(string)
		ack := claims["ack"].(string)
		if result := c.Cache.RefreshToken.Verify(jti, ack); !result {
			return nil, errors.New("refresh token verification expired")
		}
		var myToken *tokenx.Token
		if myToken, err = tokenx.Make(jwt.MapClaims{
			"jti":      jti,
			"ack":      ack,
			"username": claims["username"],
			"role":     claims["role"],
		}, time.Hour*2); err != nil {
			return nil, err
		}
		ctx.SetCookie("system_token", myToken.Value, 0, "", "", true, true)
		return myToken.Claims, nil
	})
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
