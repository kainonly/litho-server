package controller

import (
	"github.com/kainonly/iris-helper/res"
	"github.com/kainonly/iris-helper/validate"
	"log"
	"van-api/app/cache"
)

type MainController struct {
}

type LoginBody struct {
	Username string `validate:"required,min=4,max=20"`
	Password string `validate:"required,min=12,max=20"`
}

func (c *MainController) PostLogin(body LoginBody, cache *cache.Model) interface{} {
	errs := validate.Make(body, validate.Message{
		"Username": map[string]string{
			"required": "Submit missing [username] field",
		},
	})
	if errs != nil {
		return res.Error(errs)
	}
	result, err := cache.AdminGet(body.Username)
	if err != nil {
		return res.Error(err)
	}
	log.Println(result)
	return res.Ok()
}

func (c *MainController) PostVerify() interface{} {
	return res.Ok()
}

func (c *MainController) PostTest(cache *cache.Model) interface{} {
	data, err := cache.RoleGet([]string{"*"}, "resource")
	if err != nil {
		return res.Error(err)
	}
	return res.Data(data)
}
