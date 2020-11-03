package admin

import (
	curd "github.com/kainonly/iris-curd"
	"github.com/kainonly/iris-helper/hash"
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
		Originlists(model.Admin{}, body.OriginListsBody).
		OrderBy([]string{"create_time desc"}).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type listsBody struct {
	curd.ListsBody
}

func (c *Controller) PostLists(body *listsBody, mode *curd.Curd) interface{} {
	return mode.
		Lists(model.Admin{}, body.ListsBody).
		OrderBy([]string{"create_time desc"}).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type getBody struct {
	curd.GetBody
}

func (c *Controller) PostGet(body *getBody, mode *curd.Curd) interface{} {
	return mode.
		Get(model.Admin{}, body.GetBody).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type addBody struct {
	Username string `validate:"required,min=4,max=20"`
	Password string `validate:"required,min=12,max=20"`
	Role     string `validate:"required"`
	Email    string
	Phone    string
	Call     string
	Avatar   string
	Status   bool
}

func (c *Controller) PostAdd(body *addBody, mode *curd.Curd, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var err error
	var password string
	password, err = hash.Make(body.Password, hash.Option{})
	if err != nil {
		return res.Error(err)
	}
	data := model.AdminBasic{
		Username: body.Username,
		Password: password,
		Email:    body.Email,
		Phone:    body.Phone,
		Call:     body.Call,
		Avatar:   body.Avatar,
		Status:   body.Status,
	}
	return mode.
		Add().
		After(func(tx *gorm.DB) error {
			roleData := model.AdminRoleAssoc{
				Username: body.Username,
				RoleKey:  body.Role,
			}
			err = tx.Create(&roleData).Error
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
	Username string
	Password string `validate:"min=12,max=20"`
	Role     string `validate:"required_if=switch false"`
	Email    string
	Phone    string
	Call     string
	Avatar   string
	Status   bool
}

func (c *Controller) PostEdit(body *editBody, mode *curd.Curd, cache *cache.Model) interface{} {
	errs := validate.Make(body, nil)
	if errs != nil {
		return res.Error(errs)
	}
	var err error
	var password string
	if !body.Switch {
		if body.Password != "" {
			password, err = hash.Make(body.Password, hash.Option{})
			if err != nil {
				return res.Error(err)
			}
		}
	}
	data := model.AdminBasic{
		Username: body.Username,
		Password: password,
		Email:    body.Email,
		Phone:    body.Phone,
		Call:     body.Call,
		Avatar:   body.Avatar,
		Status:   body.Status,
	}
	return mode.
		Edit(model.AdminBasic{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			if !body.Switch {
				roleData := model.AdminRoleAssoc{
					Username: body.Username,
					RoleKey:  body.Role,
				}
				err = tx.Create(&roleData).Error
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
		Delete(model.AdminBasic{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(cache)
			return nil
		}).
		Exec()
}

func clearcache(cache *cache.Model) {
	cache.AdminClear()
}
