package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/weplanx/support/api"
	"laboratory/common"
)

type Controller struct {
	*InjectController
	*api.API
}

type InjectController struct {
	common.App
	Service *Service
}

func (x *Controller) Create(c *gin.Context) interface{} {
	return x.API.Create(c)
}
