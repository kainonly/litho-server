package acl

import (
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
	"github.com/kainonly/gin-extra/datatype"
	"gorm.io/gorm"
	"lab-api/application/common"
	"lab-api/application/model"
	"strings"
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
		curd.Plan(model.Acl{}, body),
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
		curd.Plan(model.Acl{}, body),
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
		curd.Plan(model.Acl{}, body),
	).Get()
}

type addBody struct {
	Key    string              `binding:"required"`
	Name   datatype.JSONObject `binding:"required"`
	Read   []string
	Write  []string
	Status bool
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Acl{
		Key:    body.Key,
		Name:   body.Name,
		Read:   strings.Join(body.Read, ","),
		Write:  strings.Join(body.Write, ","),
		Status: body.Status,
	}
	return c.Curd.Operates(
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
	).Add(&data)
}

type editBody struct {
	curd.Edit
	Key    string              `binding:"switch"`
	Name   datatype.JSONObject `binding:"switch"`
	Read   []string
	Write  []string
	Status bool `binding:"switch"`
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Acl{
		Key:    body.Key,
		Name:   body.Name,
		Read:   strings.Join(body.Read, ","),
		Write:  strings.Join(body.Write, ","),
		Status: body.Status,
	}
	return c.Curd.Operates(
		curd.Plan(model.Acl{}, body),
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
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
		curd.Plan(model.Acl{}, body),
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
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
	c.Db.Model(&model.Acl{}).
		Where("`key` = ?", body.Key).
		Count(&count)

	return gin.H{
		"exists": count != 0,
	}
}

func (c *Controller) clearcache() {
	c.Redis.Acl.Clear()
	c.Redis.Role.Clear()
}
