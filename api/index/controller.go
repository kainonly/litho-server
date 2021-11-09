package index

import (
	"github.com/gin-gonic/gin"
	"laboratory/common"
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

func (x *Controller) Index(c *gin.Context) interface{} {
	return gin.H{
		"msg": "hi",
	}
}
