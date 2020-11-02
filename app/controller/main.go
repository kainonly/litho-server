package controller

import (
	"github.com/kainonly/iris-helper/res"
	"github.com/kainonly/iris-helper/validate"
	"log"
	"van-api/app/cache"
)

type Controller struct {
}

type LoginBody struct {
	Username string `validate:"required,min=4,max=20"`
	Password string `validate:"required,min=12,max=20"`
}

func (c *Controller) PostLogin(body LoginBody, cache *cache.Model) interface{} {
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

func (c *Controller) PostVerify() interface{} {
	return res.Ok()
}
