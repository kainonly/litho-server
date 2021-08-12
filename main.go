package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/go-bit"
	"github.com/kainonly/go-bit/crud"
	"log"
)

func main() {
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	config, err := bit.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	s, err := Boot(config)
	if err != nil {
		log.Fatalln(err)
	}
	app.GET("/", crud.Bind(s.Index.Index))
	resourceRoute := app.Group("resource")
	{
		resource := s.Resource
		resourceRoute.POST("originLists", crud.Bind(resource.OriginLists))
		resourceRoute.POST("lists", crud.Bind(resource.Lists))
		resourceRoute.POST("get", crud.Bind(resource.Get))
		resourceRoute.POST("add", crud.Bind(resource.Add))
		resourceRoute.POST("edit", crud.Bind(resource.Edit))
		resourceRoute.POST("delete", crud.Bind(resource.Delete))
	}
	app.Run(":8000")
}
