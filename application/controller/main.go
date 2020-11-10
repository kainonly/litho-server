package controller

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/hash"
	"github.com/kainonly/gin-extra/helper/res"
	"github.com/kainonly/gin-extra/helper/str"
	"github.com/kainonly/gin-extra/helper/token"
	"strings"
	"taste-api/application/common"
	"time"
)

type Controller struct {
}

type LoginBody struct {
	Username string `binding:"required,min=4,max=20"`
	Password string `binding:"required,min=12,max=20"`
}

func (c *Controller) Login(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body LoginBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	result, err := app.Cache.Admin.Get(body.Username)
	if err != nil {
		return res.Error(err)
	}
	ok, err := hash.Verify(body.Password, result["password"].(string))
	if err != nil {
		return res.Error(err)
	}
	if !ok {
		return res.Error("user login password is incorrect")
	}
	jti := str.Uuid()
	ack := str.Random(8)
	app.Cache.RefreshToken.TokenFactory(jti.String(), ack, time.Hour*24*7)
	myToken, err := token.Make("system", jwt.MapClaims{
		"jti":      jti,
		"ack":      ack,
		"username": result["username"],
		"role":     strings.Split(result["role"].(string), ","),
	})
	if err != nil {
		return res.Error(err)
	}
	ctx.SetCookie("system_token", myToken.String(), 0, "", "", true, true)
	return res.Ok()
}

func (c *Controller) Verify(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	tokenString, err := ctx.Cookie("system_token")
	if err != nil {
		return res.Error(err)
	}
	var claims jwt.MapClaims
	claims, err = token.Verify("system", tokenString,
		func(option token.Option) (newClaims jwt.MapClaims, err error) {
			jti := claims["jti"].(string)
			ack := claims["ack"].(string)
			result := app.Cache.RefreshToken.TokenVerify(jti, ack)
			if !result {
				err = errors.New("refresh token verification expired")
				return
			}
			var myToken *token.Token
			myToken, err = token.Make("system", jwt.MapClaims{
				"jti":      jti,
				"ack":      ack,
				"username": claims["username"],
				"role":     claims["role"],
			})
			if err != nil {
				return
			}
			newClaims = myToken.Claims()
			ctx.SetCookie("system_token", myToken.String(), 0, "", "", true, true)
			return
		},
	)
	if err != nil {
		return res.Error(err)
	}
	return res.Ok()
}

func (c *Controller) Resource(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	resource, err := app.Cache.Resource.Get()
	if err != nil {
		return res.Error(err)
	}
	//roleKeyids := []string{"*"}
	//role, err := app.Cache.RoleGet(roleKeyids, "resource")
	//if err != nil {
	//	return res.Error(err)
	//}
	return res.Data(resource)
}
