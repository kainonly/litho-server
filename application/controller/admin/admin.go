package admin

import "github.com/gin-gonic/gin"

type Controller struct {
}

func (c *Controller) Hello(ctx *gin.Context) interface{} {
	return gin.H{
		"hi": "kain",
	}
}
