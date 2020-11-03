package role

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
		Originlists(model.Role{}, body.OriginListsBody).
		OrderBy([]string{"create_time desc"}).
		Exec()
}

type listsBody struct {
	curd.ListsBody
}

func (c *Controller) PostLists(body *listsBody, mode *curd.Curd) interface{} {
	return mode.
		Lists(model.Role{}, body.ListsBody).
		OrderBy([]string{"create_time desc"}).
		Exec()
}

type getBody struct {
	curd.GetBody
}

func (c *Controller) PostGet(body *getBody, mode *curd.Curd) interface{} {
	return mode.
		Get(model.Role{}, body.GetBody).
		Exec()
}

type addBody struct {
	Keyid    string     `validate:"required"`
	Name     types.JSON `validate:"required"`
	Resource []string   `validate:"required"`
	Note     string
	Status   bool
}

func (c *Controller) PostAdd(body *addBody, mode *curd.Curd, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var err error
	data := model.RoleBasic{
		Keyid:  body.Keyid,
		Name:   body.Name,
		Note:   body.Note,
		Status: body.Status,
	}
	return mode.
		Add().
		After(func(tx *gorm.DB) error {
			var assoc []model.RoleResourceAssoc
			for _, resourceKey := range body.Resource {
				assoc = append(assoc, model.RoleResourceAssoc{
					RoleKey:     body.Keyid,
					ResourceKey: resourceKey,
				})
			}
			err = tx.Create(&assoc).Error
			if err != nil {
				return err
			}
			clearcache(cache)
			return nil
		}).
		Exec(&data)
}

type editBody struct {
	curd.EditBody
	Keyid    string     `validate:"required_if=switch false"`
	Name     types.JSON `validate:"required_if=switch false"`
	Resource []string   `validate:"required_if=switch false"`
	Note     string
	Status   bool
}

func (c *Controller) PostEdit(body *editBody, mode *curd.Curd, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var err error
	data := model.RoleBasic{
		Keyid:  body.Keyid,
		Name:   body.Name,
		Note:   body.Note,
		Status: body.Status,
	}
	return mode.
		Edit(model.Resource{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			if !body.Switch {
				err = tx.Where("role_key = ?", body.Keyid).
					Delete(model.RoleResourceAssoc{}).
					Error
				if err != nil {
					return err
				}
				var assoc []model.RoleResourceAssoc
				for _, resourceKey := range body.Resource {
					assoc = append(assoc, model.RoleResourceAssoc{
						RoleKey:     body.Keyid,
						ResourceKey: resourceKey,
					})
				}
				err = tx.Create(&assoc).Error
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

func (c *Controller) PostDelete(body *deleteBody, mode *curd.Curd, cache *cache.Model) interface{} {
	return mode.
		Delete(model.RoleBasic{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(cache)
			return nil
		}).
		Exec()
}

type validedkeyBody struct {
	Keyid string `validate:"required"`
}

func (c *Controller) PostValidedkey(body *validedkeyBody, db *gorm.DB) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var count int64
	db.Model(&model.RoleBasic{}).
		Where("keyid = ?", body.Keyid).
		Count(&count)

	return res.Data(count != 0)
}

func clearcache(cache *cache.Model) {
	cache.RoleClear()
	cache.AdminClear()
}
