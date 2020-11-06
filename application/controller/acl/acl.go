package acl

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"github.com/kainonly/gin-extra/helper/res"
	"gorm.io/gorm"
	"taste-api/application/cache"
	"taste-api/application/common"
	"taste-api/application/common/types"
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
		Originlists(model.Acl{}, body.OriginListsBody).
		OrderBy([]string{"create_time desc"}).
		Exec()
}

type listsBody struct {
	operates.ListsBody
}

func (c *Controller) Lists(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body listsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Lists(model.Acl{}, body.ListsBody).
		OrderBy([]string{"create_time desc"}).
		Exec()
}

type getBody struct {
	operates.GetBody
}

func (c *Controller) Get(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body getBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	return app.Curd.
		Get(model.Acl{}, body.GetBody).
		Exec()
}

type addBody struct {
	Keyid  string     `binding:"required"`
	Name   types.JSON `binding:"required"`
	Read   string
	Write  string
	Status bool
}

func (c *Controller) Add(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	data := model.Acl{
		Keyid:  body.Keyid,
		Name:   body.Name,
		Read:   body.Read,
		Write:  body.Write,
		Status: body.Status,
	}
	return app.Curd.
		Add().
		Exec(&data)
}

type editBody struct {
	operates.EditBody
	Keyid  string     `binding:"required_if=switch false"`
	Name   types.JSON `binding:"required_if=switch false"`
	Read   string
	Write  string
	Status bool
}

func (c *Controller) Edit(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	data := model.Acl{
		Keyid:  body.Keyid,
		Name:   body.Name,
		Read:   body.Read,
		Write:  body.Write,
		Status: body.Status,
	}
	return app.Curd.
		Edit(model.Acl{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec(data)
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
		Delete(model.Acl{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec()
}

type validedkeyBody struct {
	Keyid string `binding:"required"`
}

func (c *Controller) Validedkey(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body validedkeyBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	var count int64
	app.Db.Model(&model.Acl{}).
		Where("keyid = ?", body.Keyid).
		Count(&count)

	return res.Data(count != 0)
}

func clearcache(cache *cache.Model) {
	cache.AclClear()
	cache.RoleClear()
}
