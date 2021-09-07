package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit/support"
	"go/format"
	"io/fs"
	"io/ioutil"
)

type Index struct {
	*Dependency
}

func NewIndex(d Dependency) *Index {
	return &Index{
		Dependency: &d,
	}
}

func (x *Index) Index(c *gin.Context) interface{} {
	return x.IndexService.Version()
}

func (x *Index) Setup(c *gin.Context) interface{} {
	if err := support.GenerateResources(x.Db); err != nil {
		return err
	}
	return "ok"
}

func (x *Index) Sync(c *gin.Context) interface{} {
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
