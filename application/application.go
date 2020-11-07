package application

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/mvc"
	"taste-api/application/common"
	"taste-api/application/controller"
	"taste-api/application/controller/acl"
	"taste-api/application/controller/admin"
	"taste-api/application/controller/policy"
	"taste-api/application/controller/resource"
	"taste-api/application/controller/role"
	"taste-api/routes"
)

func Application(router *gin.Engine, dependency common.Dependency) {
	router.GET("/", routes.Default)
	system := router.Group("/system")
	{
		m := mvc.Factory(system, dependency)
		m.AutoController("/main", new(controller.Controller))
		m.AutoController("/acl", new(acl.Controller))
		m.AutoController("/resource", new(resource.Controller))
		m.AutoController("/policy", new(policy.Controller))
		m.AutoController("/role", new(role.Controller))
		m.AutoController("/admin", new(admin.Controller))
	}
}
