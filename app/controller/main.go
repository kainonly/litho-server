package controller

import (
	"github.com/kataras/iris/v12"
	"van-api/helper/res"
	"van-api/helper/validate"
)

type MainController struct {
}

type LoginBody struct {
	Username string `validate:"required,min=4,max=20"`
	Password string `validate:"required,min=12,max=20"`
}

func (c *MainController) PostLogin(ctx iris.Context) interface{} {
	var body LoginBody
	err := ctx.ReadJSON(&body)
	if err != nil {
		return res.Error(err)
	}
	errs := validate.Make(body, validate.Message{
		"Username": map[string]string{
			"required": "Submit missing [username] field",
		},
	})
	if errs != nil {
		return res.Error(errs)
	}
	return res.Ok()
}

func (c *MainController) PostVerify(ctx iris.Context) interface{} {
	return res.Ok()
}
