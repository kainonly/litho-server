package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-helper/authx"
	"github.com/kainonly/gin-helper/hash"
	"gorm.io/gorm"
	"lab-api/service"
)

type Main struct {
	admin *service.Admins
	auth  *authx.Auth
}

func NewMain(admins *service.Admins, auth *authx.Auth) *Main {
	return &Main{admins, auth}
}

func (x *Main) Index(c *gin.Context) interface{} {
	return gin.H{
		"version": 1.0,
	}
}

func (x *Main) Login(c *gin.Context) interface{} {
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
	if result, err := hash.Verify(body.Password, data.Password); err != nil || !result {
		return errors.New("the account does not exist or the password is incorrect")
	}
	if _, err := x.auth.Create(c, data.Username, data.ID, map[string]interface{}{}); err != nil {
		return err
	}
	return "ok"
}

func (x *Main) Verify(c *gin.Context) interface{} {
	if err := x.auth.Verify(c); err != nil {
		return err
	}
	return "ok"
}

func (x *Main) Logout(c *gin.Context) interface{} {
	if err := x.auth.Destory(c); err != nil {
		return err
	}
	return "ok"
}
