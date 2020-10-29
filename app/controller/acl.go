package controller

import (
	"log"
	"van-api/app/model"
	"van-api/curd"
	"van-api/helper/res"
	"van-api/helper/validate"
	"van-api/types"
)

type AclController struct {
}

type OriginListsBody struct {
	curd.OriginListsBody
}

func (c *AclController) PostOriginlists(body *OriginListsBody, mode *curd.Curd) interface{} {
	return mode.
		Originlists(model.Acl{}, body.OriginListsBody).
		Where(curd.Conditions{
			[]interface{}{"status", "=", "1"},
		}).
		Field([]string{"id", "name", "read", "write"}).
		Exec()
}

type ListsBody struct {
	curd.ListsBody
}

func (c *AclController) PostLists(body *ListsBody, mode *curd.Curd) interface{} {
	return mode.
		Lists(model.Acl{}, body.ListsBody).
		Exec()
}

type GetBody struct {
	curd.GetBody
}

func (c *AclController) PostGet(body *GetBody, mode *curd.Curd) interface{} {
	return mode.
		Get(model.Acl{}, body.GetBody).
		Field([]string{"id", "name", "read", "write"}).
		Exec()
}

type TestAdd struct {
	Keyid string     `validate:"required"`
	Name  types.JSON `validate:"required"`
	Read  string     `validate:"required"`
	Write string     `validate:"required"`
}

func (c *AclController) PostAdd(body *TestAdd, mode *curd.Curd) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	data := model.Acl{
		Keyid: body.Keyid,
		Name:  body.Name,
		Read:  body.Read,
		Write: body.Write,
	}
	return mode.Add().Exec(data)
}

type TestEdit struct {
	curd.EditBody
	Keyid string
	Name  types.JSON
	Read  string
	Write string
}

func (c *AclController) PostEdit(body *TestEdit, mode *curd.Curd) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	log.Println(body.Switch)
	data := model.Acl{
		Keyid: body.Keyid,
		Name:  body.Name,
		Read:  body.Read,
		Write: body.Write,
	}
	return mode.
		Edit(model.Acl{}, body.EditBody).
		Exec(data)
}

type TestDelete struct {
	curd.DeleteBody
}

func (c *AclController) PostDelete(body *TestDelete, mode *curd.Curd) interface{} {
	return mode.
		Delete(model.Acl{}, body.DeleteBody).
		Exec()
}
