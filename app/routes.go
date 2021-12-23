package app

import (
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) interface{} {
	return gin.H{"msg": "hi"}
}
