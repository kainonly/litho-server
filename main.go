package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit"
	"github.com/kainonly/go-bit/crud"
	"log"
)

func main() {
	config, err := bit.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	s, err := Boot(config)
	if err != nil {
		log.Fatalln(err)
	}
	app.GET("/", s.Index.Index)
	resourceRoute := app.Group("resource")
	{
		resource := s.Resource
		resourceRoute.POST("get", crud.Bind(resource.Get))

	}
	app.Run(":8000")
}
