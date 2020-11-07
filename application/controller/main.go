package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/hash"
	"github.com/kainonly/gin-extra/helper/res"
	"github.com/kainonly/gin-extra/helper/token"
	"strings"
	"taste-api/application/common"
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
	result, err := app.Cache.AdminGet(body.Username)
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
	tokenString, err := token.Make("system", jwt.MapClaims{
		"username": result["username"],
		"role":     strings.Split(result["role"].(string), ","),
	})
	if err != nil {
		return res.Error(err)
	}
	ctx.SetCookie("system_token", tokenString, 0, "", "", true, true)
	return res.Ok()
}

func (c *Controller) Verify() interface{} {
	return res.Error(nil)
}
