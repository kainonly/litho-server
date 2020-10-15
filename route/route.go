package route

import "github.com/kataras/iris/v12"

func Default(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"version": 1.0,
	})
}

func Option(ctx iris.Context) {
	ctx.StatusCode(200)
}
