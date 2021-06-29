package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-planx/mvc"
	"lab-api/controller"
)

func Initialize(
	route *gin.Engine,
	main *controller.Main,
) {
	routes := [][]interface{}{
		{"GET", "/", main.Index},
	}

	for _, r := range routes {
		handlers := []gin.HandlerFunc{mvc.Bind(r[2])}
		for _, ext := range r[3:] {
			handlers = append(handlers, ext.(gin.HandlerFunc))
		}
		route.Handle(r[0].(string), r[1].(string), handlers...)
	}
}
