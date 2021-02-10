package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
	"github.com/kainonly/gin-extra/datatype"
	"gorm.io/gorm"
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
		curd.Plan(model.Resource{}, body),
		curd.OrderBy(curd.Orders{"sort": "asc"}),
	).Originlists()
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
		curd.Plan(model.Resource{}, body),
	).Get()
}

type addBody struct {
	Key    string `binding:"required"`
	Parent string
	Name   datatype.JSONObject `binding:"required"`
	Nav    bool
	Router bool
	Policy bool
	Icon   string
	Sort   uint8
	Status bool
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Resource{
		Key:    body.Key,
		Parent: body.Parent,
		Name:   body.Name,
		Nav:    body.Nav,
		Router: body.Router,
		Policy: body.Policy,
		Icon:   body.Icon,
		Sort:   body.Sort,
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
	Key    string
	Parent string
	Name   datatype.JSONObject
	Nav    bool
	Router bool
	Policy bool
	Icon   string
	Sort   uint8
	Status bool
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	var prevData model.Resource
	if !body.Switch {
		c.Db.Where("id = ?", body.Id).
			Find(&prevData)
	}
	data := model.Resource{
		Key:    body.Key,
		Parent: body.Parent,
		Name:   body.Name,
		Nav:    body.Nav,
		Router: body.Router,
		Policy: body.Policy,
		Icon:   body.Icon,
		Sort:   body.Sort,
		Status: body.Status,
	}
	return c.Curd.Operates(
		curd.Plan(model.Resource{}, body),
		curd.After(func(tx *gorm.DB) error {
			if !body.Switch && body.Key != prevData.Key {
				if err = tx.Model(&model.Resource{}).
					Where("parent = ?", body.Key).
					Update("parent", body.Key).
					Error; err != nil {
					return err
				}
			}
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
	var data model.Resource
	c.Db.Where("id in ?", body.Id).First(&data)
	var count int64
	c.Db.Model(&model.Resource{}).Where("parent = ?", data.Key).Count(&count)
	if count != 0 {
		return errors.New("A subset of nodes cannot be deleted")
	}
	return c.Curd.Operates(
		curd.Plan(model.Resource{}, body),
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
	).Delete()
}

type sortBody struct {
	Data []map[string]interface{} `binding:"required"`
}

func (c *Controller) Sort(ctx *gin.Context) interface{} {
	var body sortBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	if err = c.Db.Transaction(func(tx *gorm.DB) error {
		for _, value := range body.Data {
			if err = tx.Model(&model.Resource{}).
				Where("id = ?", value["id"]).
				Update("sort", value["sort"]).
				Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return true

}

type bindingkeyBody struct {
	Key string `binding:"required"`
}

func (c *Controller) Bindingkey(ctx *gin.Context) interface{} {
	var body bindingkeyBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	var count int64
	c.Db.Model(&model.Resource{}).
		Where("`key` = ?", body.Key).
		Count(&count)

	return gin.H{
		"exists": count != 0,
	}
}

func (c *Controller) clearcache() {
	c.Cache.Resource.Clear()
	c.Cache.Role.Clear()
}
