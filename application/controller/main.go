package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/hash"
	"github.com/kainonly/gin-extra/helper/res"
	"github.com/kainonly/gin-extra/helper/str"
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
	jti := str.Uuid()
	//ack := str.Random(8)

	tokenString, err := token.Make("system", jwt.MapClaims{
		"jti":      jti,
		"username": result["username"],
		"role":     strings.Split(result["role"].(string), ","),
	})
	if err != nil {
		return res.Error(err)
	}
	ctx.SetCookie("system_token", tokenString, 0, "", "", true, true)
	return res.Ok()
}

func (c *Controller) Verify(ctx *gin.Context) interface{} {
	tokenString, err := ctx.Cookie("system_token")
	if err != nil {
		return res.Error(err)
	}
	_, err = token.Verify("system", tokenString, func(option token.Option) (claims jwt.MapClaims, err error) {
		return
	})
	if err != nil {
		return res.Error(err)
	}
	return res.Ok()
}

func (c *Controller) Resource(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	resource, err := app.Cache.ResourceGet()
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
