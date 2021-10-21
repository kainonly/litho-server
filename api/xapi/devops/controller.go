package devops

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/support/basic"
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
	if err := basic.GenerateSchema(c, x.Db); err != nil {
		return err
	}
	if err := basic.GeneratePage(c, x.Db); err != nil {
		return err
	}
	if err := x.Service.InitData(c); err != nil {
		return err
	}
	return "ok"
}
