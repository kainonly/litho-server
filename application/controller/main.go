package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/hash"
	"github.com/kainonly/gin-extra/helper/str"
	"github.com/kainonly/gin-extra/helper/token"
	"lab-api/application/common"
	"time"
)

type Controller struct {
	*common.Dependency
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
	jti := str.Uuid()
	ack := str.Random(8)
	c.Cache.RefreshToken.TokenFactory(jti.String(), ack, time.Hour*24*7)
	stoken, _ := token.Make(jwt.MapClaims{
		"jti":      jti,
		"iss":      "api",
		"aud":      []string{"system"},
		"ack":      ack,
		"username": admin["username"],
		"role":     admin["role"],
	}, time.Hour*2)
	ctx.SetCookie("system_token", stoken.String(), 0, "", "", true, true)
	return true
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	tokenString, _ := ctx.Cookie("system_token")
	token.Verify(tokenString, func(claims jwt.MapClaims) (jwt.MapClaims, error) {
		var err error
		jti := claims["jti"].(string)
		ack := claims["ack"].(string)
		if result := c.Cache.RefreshToken.TokenVerify(jti, ack); !result {
			return nil, errors.New("refresh token verification expired")
		}
		var myToken *token.Token
		if myToken, err = token.Make(jwt.MapClaims{
			"jti":      jti,
			"ack":      ack,
			"username": claims["username"],
			"role":     claims["role"],
		}, time.Hour*2); err != nil {
			return nil, err
		}
		ctx.SetCookie("system_token", myToken.String(), 0, "", "", true, true)
		return myToken.Claims(), nil
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
