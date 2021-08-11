package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	s, err := Boot()
	if err != nil {
		log.Fatalln(err)
	}
	r.GET("/", s.Index.Index)

	r.Run(":8000")
}
