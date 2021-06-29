package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-bit"
	"lab-api/controller"
)

func Initialize(
	route *gin.Engine,
	main *controller.Main,
	acl *controller.Acl,
) {
	route.GET("/", bit.Bind(main.Index))
	ACL := route.Group("/acl")
	{
		ACL.POST("/get", bit.Bind(acl.Get))

	}
}
