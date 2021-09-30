package devops

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/support/basic"
	"go/format"
	"io/fs"
	"io/ioutil"
	"laboratory/common"
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
	if err := basic.GenerateSchema(tx); err != nil {
		return err
	}
	if err := basic.GeneratePage(tx); err != nil {
		return err
	}
	if err, ok := x.Sync(c).(error); ok {
		return err
	}
	return "ok"
}

func (x *Controller) Sync(c *gin.Context) interface{} {
	buf, err := basic.GenerateModel(x.Db.WithContext(c))
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
