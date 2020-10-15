package controller

import (
	"github.com/kataras/iris/v12"
	"log"
	"van-api/app/utils/res"
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
	return res.Result(iris.Map{
		"name": "kain",
	})
}
