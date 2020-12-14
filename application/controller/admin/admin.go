package admin

import (
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
	"github.com/kainonly/gin-extra/helper/hash"
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
		curd.Plan(model.Admin{}, body),
		curd.Field([]string{"password"}, true),
	).Originlists()
}

type listsBody struct {
	curd.Lists
}

func (c *Controller) Lists(ctx *gin.Context) interface{} {
	var body listsBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.Operates(
		curd.Plan(model.Admin{}, body),
		curd.Field([]string{"password"}, true),
	).Lists()
}

type getBody struct {
	curd.Get
}

func (c *Controller) Get(ctx *gin.Context) interface{} {
	var body getBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	return c.Curd.Operates(
		curd.Plan(model.Admin{}, body),
		curd.Field([]string{"password"}, true),
	).Get()
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
	return c.Curd.Operates(
		curd.After(func(tx *gorm.DB) error {
			roleData := model.AdminRoleRel{
				Username: body.Username,
				RoleKey:  body.Role,
			}
			if err = tx.Create(&roleData).Error; err != nil {
				return err
			}
			c.clearcache()
			return nil
		}),
	).Add(&data)
}

type editBody struct {
	curd.Edit
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
	return c.Curd.Operates(
		curd.After(func(tx *gorm.DB) error {
			if !body.Switch {
				roleData := model.AdminRoleRel{
					Username: body.Username,
					RoleKey:  body.Role,
				}
				if err = tx.Create(&roleData).Error; err != nil {
					return err
				}
			}
			c.clearcache()
			return nil
		}),
	).Edit(data)
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
		curd.Plan(model.AdminBasic{}, body),
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
	).Delete()
}

func (c *Controller) clearcache() {
	c.Cache.Admin.Clear()
}
