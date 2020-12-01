package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-curd/operates"
	"github.com/kainonly/gin-curd/typ"
	"github.com/kainonly/gin-extra/helper/hash"
	"gorm.io/gorm"
	"lab-api/application/cache"
	"lab-api/application/common"
	"lab-api/application/model"
)

type Controller struct {
	*common.Dependency
}

type originListsBody struct {
	operates.OriginListsBody
}

func (c *Controller) OriginLists(ctx *gin.Context) interface{} {
	var body originListsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Originlists(model.Admin{}, body.OriginListsBody).
		OrderBy(typ.Orders{"create_time": "desc"}).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type listsBody struct {
	operates.ListsBody
}

func (c *Controller) Lists(ctx *gin.Context) interface{} {
	var body listsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Lists(model.Admin{}, body.ListsBody).
		OrderBy(typ.Orders{"create_time": "desc"}).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type getBody struct {
	operates.GetBody
}

func (c *Controller) Get(ctx *gin.Context) interface{} {
	var body getBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Get(model.Admin{}, body.GetBody).
		Field([]string{"id", "username", "role", "call", "email", "phone", "avatar", "status"}).
		Exec()
}

type addBody struct {
	Username string `binding:"required,min=4,max=20"`
	Password string `binding:"required,min=12,max=20"`
	Role     string `binding:"required"`
	Email    string
	Phone    string
	Call     string
	Avatar   string
	Status   bool
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	password, _ := hash.Make(body.Password, hash.Option{})
	data := model.AdminBasic{
		Username: body.Username,
		Password: password,
		Email:    body.Email,
		Phone:    body.Phone,
		Call:     body.Call,
		Avatar:   body.Avatar,
		Status:   body.Status,
	}
	return c.Curd.
		Add().
		After(func(tx *gorm.DB) error {
			roleData := model.AdminRoleRel{
				Username: body.Username,
				RoleKey:  body.Role,
			}
			if err = tx.Create(&roleData).Error; err != nil {
				return err
			}
			clearcache(c.Cache)
			return nil
		}).
		Exec(&data)
}

type editBody struct {
	operates.EditBody
	Username string
	Password string `binding:"min=12,max=20"`
	Role     string `binding:"required_if=switch false"`
	Email    string
	Phone    string
	Call     string
	Avatar   string
	Status   bool
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	var password string
	if !body.Switch {
		if body.Password != "" {
			password, _ = hash.Make(body.Password, hash.Option{})
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
	return c.Curd.
		Edit(model.AdminBasic{}, body.EditBody).
		After(func(tx *gorm.DB) error {
			if !body.Switch {
				roleData := model.AdminRoleRel{
					Username: body.Username,
					RoleKey:  body.Role,
				}
				if err = tx.Create(&roleData).Error; err != nil {
					return err
				}
			}
			clearcache(c.Cache)
			return nil
		}).
		Exec(data)
}

type deleteBody struct {
	operates.DeleteBody
}

func (c *Controller) Delete(ctx *gin.Context) interface{} {
	var body deleteBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.
		Delete(model.AdminBasic{}, body.DeleteBody).
		After(func(tx *gorm.DB) error {
			clearcache(c.Cache)
			return nil
		}).
		Exec()
}

func clearcache(cache *cache.Cache) {
	cache.Admin.Clear()
}
