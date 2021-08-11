package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config, err := LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	s, err := Boot(config)
	if err != nil {
		log.Fatalln(err)
	}
	r.GET("/", s.Index.Index)
	r.Run(":8000")
}
