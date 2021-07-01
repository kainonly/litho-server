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
		rAcl.POST("/originLists", bit.Bind(acl.OriginLists))
		rAcl.POST("/lists", bit.Bind(acl.Lists))
		rAcl.POST("/get", bit.Bind(acl.Get))
		rAcl.POST("/add", bit.Bind(acl.Add))
	}
}
