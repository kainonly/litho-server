package index

import "github.com/gin-gonic/gin"

func (x *Controller) Index(c *gin.Context) interface{} {
	return x.Service.Version()
}
