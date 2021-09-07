package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/support"
	"go/format"
	"io/fs"
	"io/ioutil"
	"lab-api/model"
)

type Developer struct {
	*Dependency
}

func NewDeveloper(d Dependency) *Developer {
	return &Developer{
		Dependency: &d,
	}
}

func (x *Developer) Setup(c *gin.Context) interface{} {
	if err := support.GenerateResources(x.Db); err != nil {
		return err
	}
	return "ok"
}

func (x *Developer) Sync(c *gin.Context) interface{} {
	buf, err := support.GenerateModels(x.Db)
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

func (x *Developer) Migrate(c *gin.Context) interface{} {
	var body struct {
		Key string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		return err
	}
	model.AutoMigrate(x.Db, body.Key)
	return "ok"
}
