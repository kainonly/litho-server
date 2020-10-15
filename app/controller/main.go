package controller

import (
	"github.com/kataras/iris/v12"
)

type MainController struct {
}

func (c *MainController) Get() interface{} {
	return iris.Map{
		"name": "kain",
	}
}
