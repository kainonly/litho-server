package index

import (
	"api/common"
	"errors"
	"github.com/gin-gonic/gin"
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

func (x *Controller) Ms(c *gin.Context) interface{} {
	return errors.New("this is a test")
}
