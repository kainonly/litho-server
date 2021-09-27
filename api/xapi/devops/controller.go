package devops

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/support"
	"go/format"
	"io/fs"
	"io/ioutil"
	"lab-api/common"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
}

func (x *Controller) Setup(c *gin.Context) interface{} {
	tx := x.Db.WithContext(c)
	if err := support.GenerateSchema(tx); err != nil {
		return err
	}
	if err := support.GeneratePage(tx); err != nil {
		return err
	}
	if err, ok := x.Sync(c).(error); ok {
		return err
	}
	return "ok"
}

func (x *Controller) Sync(c *gin.Context) interface{} {
	buf, err := support.GenerateModel(x.Db.WithContext(c))
	if err != nil {
		return err
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("./model/model_gen.go", b, fs.ModePerm); err != nil {
		return err
	}
	return "ok"
}
