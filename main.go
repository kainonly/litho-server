package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())
	s, err := Bootstrap()
	if err != nil {
		log.Fatalln(err)
	}
	index := s.Index
	route.GET("/", index.Index)
	route.Run(":8000")
}
