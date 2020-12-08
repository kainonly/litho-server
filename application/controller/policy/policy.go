package policy

import (
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
	"gorm.io/gorm"
	"lab-api/application/common"
	"lab-api/application/model"
)

type Controller struct {
	*common.Dependency
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
		curd.Plan(model.Policy{}, body),
	).Originlists()
}

type addBody struct {
	ResourceKey string `binding:"required"`
	AclKey      string `binding:"required"`
	Policy      uint8  `binding:"required"`
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	data := model.Policy{
		ResourceKey: body.ResourceKey,
		AclKey:      body.AclKey,
		Policy:      body.Policy,
	}
	return c.Curd.Operates(
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
	).Add(&data)
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
		curd.Plan(model.Policy{}, body),
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
	).Delete()
}

func (c *Controller) clearcache() {
	c.Cache.Role.Clear()
}
