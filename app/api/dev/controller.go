package dev

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/support"
	"go/format"
	"io/fs"
	"io/ioutil"
	"lab-api/model"
)

func (x *Controller) Setup(c *gin.Context) interface{} {
	tx := x.Db.WithContext(c)
	if err := support.GenerateResources(tx); err != nil {
		return err
	}
	if err, ok := x.Sync(c).(error); ok {
		return err
	}
	model.AutoMigrate(tx, "role", "admin")
	if err := support.InitSeeder(tx); err != nil {
		return err
	}
	return "ok"
}

func (x *Controller) Sync(c *gin.Context) interface{} {
	buf, err := support.GenerateModels(x.Db.WithContext(c))
	if err != nil {
		return err
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./model/model_gen.go", b, fs.ModePerm)
	if err != nil {
		return err
	}
	return "ok"
}

func (x *Controller) Migrate(c *gin.Context) interface{} {
	var body struct {
		Key []string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	model.AutoMigrate(x.Db, body.Key...)
	return "ok"
}
