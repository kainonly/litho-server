package schema

import (
	"laboratory/common"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
}

//func (x *Controller) Create(c *gin.Context) interface{} {
//	return x.API.Create(c)
//}
