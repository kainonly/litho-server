package routes

import "github.com/gin-gonic/gin"

func Default(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": 1.0,
	})
}
