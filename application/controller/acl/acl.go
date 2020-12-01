package acl

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"gorm.io/gorm"
	"lab-api/application/cache"
	"lab-api/application/common"
	"lab-api/application/common/datatype"
	"lab-api/application/model"
)

type Controller struct {
	*common.Dependency
}

type originListsBody struct {
	operates.OriginListsBody
}

func (c *Controller) OriginLists(ctx *gin.Context) interface{} {
	var body originListsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Originlists(model.Acl{}, body.OriginListsBody).
		Exec()
}

type listsBody struct {
	operates.ListsBody
}

func (c *Controller) Lists(ctx *gin.Context) interface{} {
	var body listsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Lists(model.Acl{}, body.ListsBody).
		Exec()
}

type getBody struct {
	operates.GetBody
}

func (c *Controller) Get(ctx *gin.Context) interface{} {
	var body getBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Get(model.Acl{}, body.GetBody).
		Exec()
}

type addBody struct {
	Key    string              `binding:"required"`
	Name   datatype.JSONObject `binding:"required"`
	Read   datatype.JSONArray
	Write  datatype.JSONArray
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
		Read:   body.Read,
		Write:  body.Write,
		Status: body.Status,
	}
	return c.Curd.
		Add().
		Exec(&data)
}

type editBody struct {
	operates.EditBody
	Key    string
	Name   map[string]interface{} `json:"name"`
	Read   datatype.JSONArray     `json:"read"`
	Write  datatype.JSONArray     `json:"write"`
	Status bool
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
		Read:   body.Read,
		Write:  body.Write,
		Status: body.Status,
	}
	return c.Curd.
		Edit(model.Acl{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			clearcache(c.Cache)
			return nil
		}).
		Exec(data)
}

type deleteBody struct {
	operates.DeleteBody
}

func (c *Controller) Delete(ctx *gin.Context) interface{} {
	var body deleteBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Delete(model.Acl{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(c.Cache)
			return nil
		}).
		Exec()
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
		Where("key = ?", body.Key).
		Count(&count)

	return gin.H{
		"exists": count != 0,
	}
}

func clearcache(cache *cache.Cache) {
	cache.Acl.Clear()
	cache.Role.Clear()
}
