package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"github.com/kainonly/gin-helper/res"
	"gorm.io/gorm"
	"taste-api/application/cache"
	"taste-api/application/common"
	"taste-api/application/model"
)

type Controller struct {
}

type originListsBody struct {
	operates.OriginListsBody
}

func (c *Controller) OriginLists(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body originListsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Originlists(model.Policy{}, body.OriginListsBody).
		Exec()
}

type addBody struct {
	ResourceKey string `binding:"required"`
	AclKey      string `binding:"required"`
	Policy      uint8  `binding:"required"`
}

func (c *Controller) Add(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	data := model.Policy{
		ResourceKey: body.ResourceKey,
		AclKey:      body.AclKey,
		Policy:      body.Policy,
	}
	return app.Curd.
		Add().
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec(&data)
}

type deleteBody struct {
	operates.DeleteBody
}

func (c *Controller) Delete(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body deleteBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Delete(model.Policy{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec()
}

func clearcache(cache *cache.Model) {
	cache.RoleClear()
}
