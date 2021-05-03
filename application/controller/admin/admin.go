package admin

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	curd "github.com/kainonly/gin-curd"
	"github.com/kainonly/gin-extra/hash"
	"github.com/kainonly/gin-extra/mvcx"
	"gorm.io/gorm"
	"lab-api/application/common"
	"lab-api/application/model"
	"strings"
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
		curd.Plan(model.AdminMix{}, body),
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
		curd.Plan(model.AdminMix{}, body),
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
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	result := c.Curd.Operates(
		curd.Plan(model.AdminMix{}, body),
		curd.Field([]string{"password"}, true),
	).Get()
	var count int64
	c.Db.Model(&model.Admin{}).
		Where("username = ?", auth.(jwt.MapClaims)["user"]).
		Where("status = ?", 1).
		Count(&count)
	if count != 0 {
		result.(map[string]interface{})["self"] = true
	}
	return result
}

type addBody struct {
	Username   string   `binding:"required,min=4,max=20"`
	Password   string   `binding:"required,min=12,max=20"`
	Role       []string `binding:"required"`
	Resource   []string
	Permission []string
	Email      string
	Phone      string
	Call       string
	Avatar     string
	Status     bool
}

func (c *Controller) Add(ctx *gin.Context) interface{} {
	var body addBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	password, _ := hash.Make(body.Password)
	data := model.Admin{
		Username:   body.Username,
		Password:   password,
		Permission: strings.Join(body.Permission, ","),
		Email:      body.Email,
		Phone:      body.Phone,
		Call:       body.Call,
		Avatar:     body.Avatar,
		Status:     body.Status,
	}
	return c.Curd.Operates(
		curd.After(func(tx *gorm.DB) error {
			adminRoleRels := make([]model.AdminRoleRel, len(body.Role))
			for key, val := range body.Role {
				adminRoleRels[key] = model.AdminRoleRel{
					AdminId: data.ID,
					RoleKey: val,
				}
			}
			if err = tx.Create(&adminRoleRels).Error; err != nil {
				return err
			}
			if len(body.Resource) != 0 {
				adminResourceRels := make([]model.AdminResourceRel, len(body.Resource))
				for key, val := range body.Resource {
					adminResourceRels[key] = model.AdminResourceRel{
						AdminId:     data.ID,
						ResourceKey: val,
					}
				}
				if err = tx.Create(&adminResourceRels).Error; err != nil {
					return err
				}
			}
			c.clearcache()
			return nil
		}),
	).Add(&data)
}

type editBody struct {
	curd.Edit
	Password   string   `binding:"min=12,max=20"`
	Role       []string `binding:"switch"`
	Resource   []string
	Permission []string
	Email      string
	Phone      string
	Call       string
	Avatar     string
	Status     bool `binding:"switch"`
}

func (c *Controller) Edit(ctx *gin.Context) interface{} {
	var body editBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	var count int64
	c.Db.Model(&model.Admin{}).
		Where("username = ?", auth.(jwt.MapClaims)["user"]).
		Where("status = ?", 1).
		Count(&count)
	if count != 0 {
		return mvcx.Response{
			Code: 2,
			Msg:  "Detected as currently logged in user",
		}
	}
	var password string
	if !body.Switch {
		if body.Password != "" {
			password, _ = hash.Make(body.Password)
		}
	}
	data := model.Admin{
		Password:   password,
		Permission: strings.Join(body.Permission, ","),
		Email:      body.Email,
		Phone:      body.Phone,
		Call:       body.Call,
		Avatar:     body.Avatar,
		Status:     body.Status,
	}
	return c.Curd.Operates(
		curd.After(func(tx *gorm.DB) error {
			if !body.Switch {
				if err = tx.Where("admin_id = ?", body.Id).
					Delete(&model.AdminRoleRel{}).Error; err != nil {
					return err
				}
				adminRoleRels := make([]model.AdminRoleRel, len(body.Role))
				for key, val := range body.Role {
					adminRoleRels[key] = model.AdminRoleRel{
						AdminId: body.Id.(uint64),
						RoleKey: val,
					}
				}
				if err = tx.Create(&adminRoleRels).Error; err != nil {
					return err
				}
				if len(body.Resource) != 0 {
					if err = tx.Where("admin_id = ?", body.Id).
						Delete(&model.AdminResourceRel{}).Error; err != nil {
						return err
					}
					adminResourceRels := make([]model.AdminResourceRel, len(body.Resource))
					for key, val := range body.Resource {
						adminResourceRels[key] = model.AdminResourceRel{
							AdminId:     body.Id.(uint64),
							ResourceKey: val,
						}
					}
					if err = tx.Create(&adminResourceRels).Error; err != nil {
						return err
					}
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
	auth, exists := ctx.Get("auth")
	if !exists {
		return false
	}
	var count int64
	c.Db.Model(&model.Admin{}).
		Where("id in ?", body.Id).
		Where("username = ?", auth.(jwt.MapClaims)["user"]).
		Where("status = ?", 1).
		Count(&count)
	if count != 0 {
		return mvcx.Response{
			Code: 2,
			Msg:  "Detected as currently logged in user",
		}
	}
	return c.Curd.Operates(
		curd.Plan(model.Admin{}, body),
		curd.After(func(tx *gorm.DB) error {
			c.clearcache()
			return nil
		}),
	).Delete()
}

type validedUsernameBody struct {
	Username string `binding:"required"`
}

func (c *Controller) ValidedUsername(ctx *gin.Context) interface{} {
	var body validedUsernameBody
	var err error
	if err = ctx.ShouldBindJSON(&body); err != nil {
		return err
	}
	var count int64
	c.Db.Model(&model.Admin{}).
		Where("username = ?", body.Username).
		Count(&count)

	return gin.H{
		"exists": count != 0,
	}
}

func (c *Controller) clearcache() {
	c.Redis.Admin.Clear(context.Background())
}
