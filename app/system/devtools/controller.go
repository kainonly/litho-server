package devtools

import (
	"github.com/gin-gonic/gin"
	"go/format"
	"io/fs"
	"io/ioutil"
	"lab-api/model"
)

func (x *Controller) Setup(c *gin.Context) interface{} {
	if err := x.Service.MigrateSchema(c); err != nil {
		return err
	}
	if err := x.Service.MigrateResource(c); err != nil {
		return err
	}
	if err, ok := x.Sync(c).(error); ok {
		return err
	}
	if err := model.AutoMigrate(x.Db.WithContext(c), "role", "admin"); err != nil {
		return err
	}
	if err := x.Service.Seeder(c); err != nil {
		return err
	}
	return "ok"
}

func (x *Controller) Sync(c *gin.Context) interface{} {
	buf, err := x.Service.CreateModels(c)
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
	if err := model.AutoMigrate(x.Db.WithContext(c), body.Key...); err != nil {
		return err
	}
	return "ok"
}
