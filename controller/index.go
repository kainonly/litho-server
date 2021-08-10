package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-helper/authx"
	"github.com/kainonly/gin-helper/pwd"
	"gorm.io/gorm"
	"lab-api/service"
)

type Index struct {
	admin *service.Admin
	auth  *authx.Auth
}

func NewIndex(admin *service.Admin, auth *authx.Auth) *Index {
	return &Index{admin, auth}
}

func (x *Index) Index(c *gin.Context) interface{} {
	data, err := x.admin.FindOne(func(tx *gorm.DB) *gorm.DB {
		return tx.
			Where("username = ?", "kain").
			Where("status = ?", true)
	})
	if err != nil {
		return err
	}
	return gin.H{
		"version": "1.0",
		"data":    data,
	}
}

func (x *Index) Login(c *gin.Context) interface{} {
	var body struct {
		Username string `binding:"required"`
		Password string `binding:"required,min=12,max=20"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	data, err := x.admin.FindOne(func(tx *gorm.DB) *gorm.DB {
		return tx.
			Where("username = ?", body.Username).
			Where("status = ?", true)
	})
	if err != nil {
		return err
	}
	if result, err := pwd.Verify(body.Password, data.Password); err != nil || !result {
		return errors.New("用户不存在或口令错误")
	}
	if _, err := x.auth.Create(c, data.Username, data.ID, map[string]interface{}{}); err != nil {
		return err
	}
	return "ok"
}

func (x *Index) Verify(c *gin.Context) interface{} {
	if err := x.auth.Verify(c); err != nil {
		return err
	}
	return "ok"
}

func (x *Index) Logout(c *gin.Context) interface{} {
	if err := x.auth.Destory(c); err != nil {
		return err
	}
	return "ok"
}
