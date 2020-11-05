package routes

import (
	"github.com/gin-gonic/gin"
)

func Default() interface{} {
	return gin.H{
		"version": 1.0,
	}
}
