package resource

import (
	curd "github.com/kainonly/iris-curd"
	"github.com/kainonly/iris-helper/res"
	"github.com/kainonly/iris-helper/validate"
	"gorm.io/gorm"
	"van-api/app/cache"
	"van-api/app/model"
	"van-api/types"
)

type Controller struct {
}

type originListsBody struct {
	curd.OriginListsBody
}

func (c *Controller) PostOriginlists(body *originListsBody, mode *curd.Curd) interface{} {
	return mode.
		Originlists(model.Resource{}, body.OriginListsBody).
		OrderBy([]string{"sort asc"}).
		Exec()
}

type getBody struct {
	curd.GetBody
}

func (c *Controller) PostGet(body *getBody, mode *curd.Curd) interface{} {
	return mode.
		Get(model.Resource{}, body.GetBody).
		Exec()
}

type addBody struct {
	Keyid  string `validate:"required"`
	Parent string
	Name   types.JSON `validate:"required"`
	Nav    bool
	Router bool
	Policy bool
	Icon   string
	Sort   uint8
	Status bool
}

func (c *Controller) PostAdd(body *addBody, mode *curd.Curd, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
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
	return mode.
		Add().
		After(func(tx *gorm.DB) error {
			clearcache(cache)
			return nil
		}).
		Exec(&data)
}

type editBody struct {
	curd.EditBody
	Keyid  string `validate:"required_if=switch false"`
	Parent string
	Name   types.JSON `validate:"required_if=switch false"`
	Nav    bool
	Router bool
	Policy bool
	Icon   string
	Sort   uint8
	Status bool
}

func (c *Controller) PostEdit(body *editBody, mode *curd.Curd, cache *cache.Model, db *gorm.DB) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var err error
	var prevData model.Resource
	if !body.Switch {
		db.Where("id = ?", body.Id).
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
	return mode.
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
			clearcache(cache)
			return nil
		}).
		Exec(data)
}

type deleteBody struct {
	curd.DeleteBody
}

func (c *Controller) PostDelete(body *deleteBody, mode *curd.Curd, cache *cache.Model, db *gorm.DB) interface{} {
	var data model.Resource
	db.Where("id in ?", body.Id).First(&data)
	var count int64
	db.Model(&model.Resource{}).Where("parent = ?", data.Keyid).Count(&count)
	if count != 0 {
		return res.Error("A subset of nodes cannot be deleted")
	}

	return mode.
		Delete(model.Resource{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(cache)
			return nil
		}).
		Exec()
}

type sortBody struct {
	Data []map[string]interface{} `validate:"required"`
}

func (c *Controller) PostSort(body *sortBody, db *gorm.DB, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var err error
	err = db.Transaction(func(tx *gorm.DB) error {
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

type validatekeyBody struct {
	Keyid string `validate:"required"`
}

func (c *Controller) PostValidatekey(body *validatekeyBody, db *gorm.DB) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var count int64
	db.Model(&model.Resource{}).
		Where("keyid = ?", body.Keyid).
		Count(&count)

	return res.Data(count != 0)
}

func clearcache(cache *cache.Model) {
	cache.ResourceClear()
	cache.RoleClear()
}
