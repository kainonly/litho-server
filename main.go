package main

import (
	"github.com/gin-gonic/gin"
	"lab-api/routes"
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
	routes.Initialize(route, s)
	route.Run(":8000")
}
