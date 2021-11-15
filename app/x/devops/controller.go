package devops

import (
	"api/common"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/basic"
)

type InjectController struct {
	*common.App
	Service *Service
}

type Controller struct {
	*InjectController
}

func NewController(i *InjectController) *Controller {
	return &Controller{
		InjectController: i,
	}
}

func (x *Controller) Setup(c *gin.Context) interface{} {
	if err := basic.GenerateSchema(c, x.Db); err != nil {
		return err
	}
	if err := basic.GeneratePage(c, x.Db); err != nil {
		return err
	}
	if err := basic.GenerateRoleAndAdmin(c, x.Db); err != nil {
		return err
	}
	return gin.H{"message": "ok"}
}
