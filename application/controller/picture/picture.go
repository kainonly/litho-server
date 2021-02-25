package picture

import (
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
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
		curd.Plan(model.Picture{}, body),
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
		curd.Plan(model.Picture{}, body),
	).Lists()
}

type bulkAddBody struct {
	TypeId uint64    `binding:"required"`
	Data   []addData `binding:"required"`
}

type addData struct {
	Name string `binding:"required"`
	Url  string `binding:"required"`
}

func (c *Controller) BulkAdd(ctx *gin.Context) interface{} {
	var body bulkAddBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := make([]model.Picture, len(body.Data))
	for index, value := range body.Data {
		data[index] = model.Picture{
			TypeId: body.TypeId,
			Name:   value.Name,
			Url:    value.Url,
		}
	}
	if err = c.Db.Create(data).Error; err != nil {
		return err
	}
	return true
}

type editBody struct {
	curd.Edit
	TypeId uint64 `binding:"required"`
	Name   string `binding:"required"`
	Url    string `binding:"required"`
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Picture{
		Name: body.Name,
		Url:  body.Url,
	}
	return c.Curd.Operates(
		curd.Plan(model.Picture{}, body),
	).Edit(data)
}

type bulkEditBody struct {
	TypeId uint64   `binding:"required"`
	Ids    []uint64 `binding:"required"`
}

func (c *Controller) BulkEdit(ctx *gin.Context) interface{} {
	var body bulkEditBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	if err = c.Db.Transaction(func(tx *gorm.DB) error {
		for _, id := range body.Ids {
			tx.Model(&model.Picture{}).
				Where("id = ?", id).
				Update("type_id", body.TypeId)
		}
		return nil
	}); err != nil {
		return err
	}
	return true
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
		curd.Plan(model.Picture{}, body),
	).Delete()
}

func (c *Controller) Count(ctx *gin.Context) interface{} {
	var total int64
	tx := c.Db.Model(&model.Picture{})
	tx.Count(&total)
	values := make([]map[string]interface{}, 0)
	tx.Group("type_id").
		Select([]string{"type_id", "count(*) as size"}).
		Find(&values)
	return gin.H{
		"total":  total,
		"values": values,
	}
}
