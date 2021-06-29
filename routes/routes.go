package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-bit"
	"lab-api/controller"
)

func Initialize(
	r *gin.Engine,
	main *controller.Main,
	acl *controller.Acl,
) {
	r.GET("/", bit.Bind(main.Index))
	rAcl := r.Group("/acl")
	{
		rAcl.POST("/get", bit.Bind(acl.Get))
	}
}
