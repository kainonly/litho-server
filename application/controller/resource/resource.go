package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"github.com/kainonly/gin-helper/res"
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
		Originlists(model.Resource{}, body.OriginListsBody).
		OrderBy([]string{"sort asc"}).
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
		Get(model.Resource{}, body.GetBody).
		Exec()
}

type addBody struct {
	Keyid  string `binding:"required"`
	Parent string
	Name   types.JSON `binding:"required"`
	Nav    bool
	Router bool
	Policy bool
	Icon   string
	Sort   uint8
	Status bool
}

func (c *Controller) Add(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	data := model.Resource{
		Keyid:  body.Keyid,
		Parent: body.Parent,
		Name:   body.Name,
		Nav:    body.Nav,
		Router: body.Router,
		Policy: body.Policy,
		Icon:   body.Icon,
		Sort:   body.Sort,
		Status: body.Status,
	}
	return app.Curd.
		Add().
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec(&data)
}

type editBody struct {
	operates.EditBody
	Keyid  string `binding:"required_if=switch false"`
	Parent string
	Name   types.JSON `binding:"required_if=switch false"`
	Nav    bool
	Router bool
	Policy bool
	Icon   string
	Sort   uint8
	Status bool
}

func (c *Controller) Edit(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	var prevData model.Resource
	if !body.Switch {
		app.Db.Where("id = ?", body.Id).
			Find(&prevData)
	}
	data := model.Resource{
		Keyid:  body.Keyid,
		Parent: body.Parent,
		Name:   body.Name,
		Nav:    body.Nav,
		Router: body.Router,
		Policy: body.Policy,
		Icon:   body.Icon,
		Sort:   body.Sort,
		Status: body.Status,
	}
	return app.Curd.
		Edit(model.Resource{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			if !body.Switch && body.Keyid != prevData.Keyid {
				err = tx.Model(&model.Resource{}).
					Where("parent = ?", body.Keyid).
					Update("parent", body.Keyid).
					Error
				if err != nil {
					return err
				}
			}
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
	var data model.Resource
	app.Db.Where("id in ?", body.Id).First(&data)
	var count int64
	app.Db.Model(&model.Resource{}).Where("parent = ?", data.Keyid).Count(&count)
	if count != 0 {
		return res.Error("A subset of nodes cannot be deleted")
	}

	return app.Curd.
		Delete(model.Resource{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(app.Cache)
			return nil
		}).
		Exec()
}

type sortBody struct {
	Data []map[string]interface{} `binding:"required"`
}

func (c *Controller) Sort(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body sortBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	err = app.Db.Transaction(func(tx *gorm.DB) error {
		for _, value := range body.Data {
			err = tx.Model(&model.Resource{}).
				Where("id = ?", value["id"]).
				Update("sort", value["sort"]).
				Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return res.Ok()
	} else {
		return res.Error(err)
	}
}

type bindingkeyBody struct {
	Keyid string `binding:"required"`
}

func (c *Controller) Bindingkey(ctx *gin.Context, i interface{}) interface{} {
	app := common.Inject(i)
	var body bindingkeyBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return res.Error(err)
	}
	var count int64
	app.Db.Model(&model.Resource{}).
		Where("keyid = ?", body.Keyid).
		Count(&count)

	return res.Data(count != 0)
}

func clearcache(cache *cache.Model) {
	cache.ResourceClear()
	cache.RoleClear()
}
