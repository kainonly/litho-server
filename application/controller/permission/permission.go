package permission

import (
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
	"github.com/kainonly/gin-extra/datatype"
	"lab-api/application/common"
	"lab-api/application/model"
)

type Controller struct {
	common.Dependency
}

type originListsBody struct {
	curd.OriginLists
}

func (c *Controller) OriginLists(ctx *gin.Context) interface{} {
	var body originListsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.Operates(
		curd.Plan(model.Permission{}, body),
	).Originlists()
}

type listsBody struct {
	curd.Lists
}

func (c *Controller) Lists(ctx *gin.Context) interface{} {
	var body listsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.Operates(
		curd.Plan(model.Permission{}, body),
	).Lists()
}

type getBody struct {
	curd.Get
}

func (c *Controller) Get(ctx *gin.Context) interface{} {
	var body getBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.Operates(
		curd.Plan(model.Permission{}, body),
	).Get()
}

type addBody struct {
	Key    string              `binding:"required"`
	Name   datatype.JSONObject `binding:"required"`
	Note   string
	Status bool
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Permission{
		Key:    body.Key,
		Name:   body.Name,
		Note:   body.Note,
		Status: body.Status,
	}
	return c.Curd.Operates().Add(&data)
}

type editBody struct {
	curd.Edit
	Key    string
	Name   datatype.JSONObject
	Note   string
	Status bool
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Permission{
		Key:    body.Key,
		Name:   body.Name,
		Note:   body.Note,
		Status: body.Status,
	}
	return c.Curd.Operates(
		curd.Plan(model.Permission{}, body),
	).Edit(data)
}

type deleteBody struct {
	curd.Delete
}

func (c *Controller) Delete(ctx *gin.Context) interface{} {
	var body deleteBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.Operates(
		curd.Plan(model.Permission{}, body),
	).Delete()
}

type validedkeyBody struct {
	Key string `binding:"required"`
}

func (c *Controller) ValidedKey(ctx *gin.Context) interface{} {
	var body validedkeyBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	var count int64
	c.Db.Model(&model.Permission{}).
		Where("`key` = ?", body.Key).
		Count(&count)

	return gin.H{
		"exists": count != 0,
	}
}
