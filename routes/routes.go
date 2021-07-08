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

	system := r.Group("/system")
	{
		system.POST("/main/login", bit.Bind(main.Login))
		system.POST("/main/verify", bit.Bind(main.Verify))
		system.POST("/main/logout", bit.Bind(main.Logout))

		rAcl := system.Group("/acl")
		{
			rAcl.POST("/originLists", bit.Bind(acl.OriginLists))
			rAcl.POST("/lists", bit.Bind(acl.Lists))
			rAcl.POST("/get", bit.Bind(acl.Get))
			rAcl.POST("/add", bit.Bind(acl.Add))
			rAcl.POST("/edit", bit.Bind(acl.Edit))
			rAcl.POST("/delete", bit.Bind(acl.Delete))
		}
	}
}
