package index

import (
	"github.com/gin-gonic/gin"
	"lab-api/common"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
}

func (x *Controller) Index(c *gin.Context) interface{} {
	return "123"
}
