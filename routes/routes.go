package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-bit"
	"lab-api/controller"
)

func Initialize(
	r *gin.Engine,
	main *controller.Main,
) {
	r.GET("/", bit.Bind(main.Index))
}
