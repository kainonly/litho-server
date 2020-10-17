package controller

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"log"
	"time"
	"van-api/helper/res"
)

type MainController struct {
}

func (c *MainController) PostVerify(ctx iris.Context) interface{} {
	var data map[string]interface{}
	err := ctx.ReadJSON(&data)
	if err != nil {
		return res.Error(err.Error())
	}
	log.Println(data)
	ctx.SetCookieKV("xxx", "sdsd")
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"time": time.Now(),
	})
	tokenString, _ := token.SignedString([]byte("My Secret"))
	return res.Result(iris.Map{
		"token": tokenString,
	})
}
