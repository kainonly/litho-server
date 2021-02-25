package picture_type

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
		curd.Plan(model.PictureType{}, body),
	).Originlists()
}

type addBody struct {
	Name string `binding:"required"`
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.PictureType{
		Name: body.Name,
	}
	return c.Curd.Operates().Add(&data)
}

type editBody struct {
	curd.Edit
	Name string `binding:"required"`
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.PictureType{
		Name: body.Name,
	}
	return c.Curd.Operates(
		curd.Plan(model.PictureType{}, body),
	).Edit(data)
}

type sortBody struct {
	Data []sortData `binding:"required"`
}

type sortData struct {
	Id   uint64 `binding:"required"`
	Sort uint8  `binding:"required"`
}

func (c *Controller) Sort(ctx *gin.Context) interface{} {
	var body sortBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	if err = c.Db.Transaction(func(tx *gorm.DB) error {
		for _, data := range body.Data {
			tx.Model(&model.PictureType{}).
				Where("id = ?", data.Id).
				Update("sort", data.Sort)
		}
		return nil
	}); err != nil {
		return err
	}
	return true
}
