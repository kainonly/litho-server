package index

import (
	"api/common"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
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
