package controller

import "github.com/kataras/iris/v12"

func (c *controller) Default(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"version": 1.0,
	})
}

func (c *controller) Option(ctx iris.Context) {
	ctx.StatusCode(200)
}
