package resource

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"github.com/kainonly/gin-curd/typ"
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
		Originlists(model.Resource{}, body.OriginListsBody).
		OrderBy(typ.Orders{"sort": "asc"}).
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
		Get(model.Resource{}, body.GetBody).
		Exec()
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
	return c.Curd.
		Add().
		After(func(tx *gorm.DB) error {
			clearcache(c.Cache)
			return nil
		}).
		Exec(&data)
}

type editBody struct {
	operates.EditBody
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
	return c.Curd.
		Edit(model.Resource{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			if !body.Switch && body.Key != prevData.Key {
				if err = tx.Model(&model.Resource{}).
					Where("parent = ?", body.Key).
					Update("parent", body.Key).
					Error; err != nil {
					return err
				}
			}
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
	var data model.Resource
	c.Db.Where("id in ?", body.Id).First(&data)
	var count int64
	c.Db.Model(&model.Resource{}).Where("parent = ?", data.Key).Count(&count)
	if count != 0 {
		return errors.New("A subset of nodes cannot be deleted")
	}

	return c.Curd.
		Delete(model.Resource{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(c.Cache)
			return nil
		}).
		Exec()
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
		Where("keyid = ?", body.Key).
		Count(&count)

	return gin.H{
		"exists": count != 0,
	}
}

func clearcache(cache *cache.Cache) {
	cache.Resource.Clear()
	cache.Role.Clear()
}
