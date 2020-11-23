package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/hash"
	"github.com/kainonly/gin-extra/helper/res"
	"github.com/kainonly/gin-extra/helper/str"
	"github.com/kainonly/gin-extra/helper/token"
	"log"
	"strings"
	"taste-api/application/common"
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
	var err error
	var body LoginBody
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	var result map[string]interface{}
	if result, err = c.Cache.Admin.Get(body.Username); err != nil {
		return res.Error(err)
	}
	var ok bool
	if ok, err = hash.Verify(body.Password, result["password"].(string)); err != nil || !ok {
		return res.Error("user login password is incorrect")
	}
	jti := str.Uuid()
	ack := str.Random(8)
	c.Cache.RefreshToken.TokenFactory(jti.String(), ack, time.Hour*24*7)
	var stoken *token.Token
	if stoken, err = token.Make(jwt.MapClaims{
		"jti":      jti,
		"iss":      "api",
		"aud":      []string{"system"},
		"ack":      ack,
		"username": result["username"],
		"role":     strings.Split(result["role"].(string), ","),
	}, time.Hour*2); err != nil {
		return res.Error(err)
	}
	ctx.SetCookie("system_token", stoken.String(), 0, "", "", true, true)
	return res.Ok()
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	var err error
	var tokenString string
	if tokenString, err = ctx.Cookie("system_token"); err != nil {
		return res.Error(err)
	}
	if _, err = token.Verify(tokenString,
		func(claims jwt.MapClaims) (jwt.MapClaims, error) {
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
		},
	); err != nil {
		return res.Error(err)
	}
	return res.Ok()
}

func (c *Controller) Resource(ctx *gin.Context) interface{} {
	resource, err := c.Cache.Resource.Get()
	if err != nil {
		return res.Error(err)
	}
	var result []string
	if result, err = c.Cache.Role.Get([]string{"*"}, "resource"); err != nil {
		return err
	}
	log.Println(result)

	//roleKeyids := []string{"*"}
	//role, err := c.Cache.RoleGet(roleKeyids, "resource")
	//if err != nil {
	//	return res.Error(err)
	//}
	return res.Data(resource)
}
