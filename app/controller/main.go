package controller

import "github.com/kataras/iris/v12"

func (c *controller) MainVerify(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"error": 1,
		"msg":   "fail",
	})
}
