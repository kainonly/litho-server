package policy

import (
	curd "github.com/kainonly/iris-curd"
	"github.com/kainonly/iris-helper/res"
	"github.com/kainonly/iris-helper/validate"
	"gorm.io/gorm"
	"van-api/app/cache"
	"van-api/app/model"
)

type Controller struct {
}

type originListsBody struct {
	curd.OriginListsBody
}

func (c *Controller) PostOriginlists(body *originListsBody, mode *curd.Curd) interface{} {
	return mode.
		Originlists(model.Policy{}, body.OriginListsBody).
		Exec()
}

type addBody struct {
	ResourceKey string `validate:"required"`
	AclKey      string `validate:"required"`
	Policy      uint8  `validate:"required"`
}

func (c *Controller) PostAdd(body *addBody, mode *curd.Curd, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	data := model.Policy{
		ResourceKey: body.ResourceKey,
		AclKey:      body.AclKey,
		Policy:      body.Policy,
	}
	return mode.
		Add().
		After(func(tx *gorm.DB) error {
			clearcache(cache)
			return nil
		}).
		Exec(&data)
}

type deleteBody struct {
	curd.DeleteBody
}

func (c *Controller) PostDelete(body *deleteBody, mode *curd.Curd, cache *cache.Model) interface{} {
	return mode.
		Delete(model.Policy{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(cache)
			return nil
		}).
		Exec()
}

func clearcache(cache *cache.Model) {
	cache.RoleClear()
}
