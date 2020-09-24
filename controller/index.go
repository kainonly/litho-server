package controller

import "github.com/kataras/iris/v12"

func (c *controller) Index(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"version": 1.0,
	})
}
