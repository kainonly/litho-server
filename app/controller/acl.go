package controller

import (
	"errors"
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/gorm"
	"van-api/app/model"
	"van-api/curd"
	"van-api/helper/res"
	"van-api/helper/validate"
	"van-api/types"
)

type AclController struct {
}

func (c *AclController) BeforeActivation(b mvc.BeforeActivation) {
}

type TestRead struct {
	curd.Body
}

func (c *AclController) PostOriginlists(body *TestRead, mode *curd.Curd) interface{} {
	return mode.
		Originlists(model.Acl{}, body).
		Where(curd.Conditions{
			[]interface{}{"status", "=", "1"},
		}).
		Field([]string{"id", "name", "read", "write"}).
		Result()
}

func (c *AclController) PostLists(body *TestRead, mode *curd.Curd) interface{} {
	return mode.
		Lists(model.Acl{}, body).
		Result()
}

func (c *AclController) PostGet(body *TestRead, mode *curd.Curd) interface{} {
	return mode.
		Get(model.Acl{}, body).
		Query(func(tx *gorm.DB) {
			tx.Where("id = ?", "2")
		}).
		Field([]string{"id", "name", "read", "write"}).
		Result()
}

type TestAdd struct {
	Key   string     `validate:"required"`
	Name  types.JSON `validate:"required"`
	Read  string     `validate:"required"`
	Write string     `validate:"required"`
}

func (c *AclController) PostAdd(body *TestAdd, mode *curd.Curd) interface{} {
	errs := validate.Make(body, validate.Message{
		"Username": map[string]string{
			"required": "Submit missing [username] field",
		},
	})
	if errs != nil {
		return res.Error(errs)
	}
	data := model.Acl{
		Key:   body.Key,
		Name:  body.Name,
		Read:  body.Read,
		Write: body.Write,
	}
	return mode.
		Add(&data, body).
		After(func(tx *gorm.DB) error {
			return errors.New("test error")
		}).
		Result()
}
