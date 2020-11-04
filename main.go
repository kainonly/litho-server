package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Body struct {
	Name string `binding:"required"`
}

func main() {
	r := gin.Default()
	r.POST("/test", func(c *gin.Context) {
		var err error
		var body Body
		if err = c.ShouldBindJSON(&body); err != nil {
			c.JSON(200, gin.H{
				"error": 1,
				"msg":   err.Error(),
			})
			return
		}
		log.Println(body.Name)
		c.JSON(200, gin.H{
			"error": 0,
		})
		return
	})
	r.Run()
}
