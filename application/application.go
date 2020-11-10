package application

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/token"
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
	cfg := dependency.Config
	token.Key = []byte(cfg.App.Key)
	token.Options = cfg.Token
	router.GET("/", routes.Default)
	system := router.Group("/system")
	{
		m := mvc.Factory(system, &dependency)
		m.AutoController(mvc.Auto{
			Path:       "/main",
			Controller: new(controller.Controller),
		})
		m.AutoController(mvc.Auto{
			Path:       "/acl",
			Controller: new(acl.Controller),
		})
		m.AutoController(mvc.Auto{
			Path:       "/resource",
			Controller: new(resource.Controller),
		})
		m.AutoController(mvc.Auto{
			Path:       "/policy",
			Controller: new(policy.Controller),
		})
		m.AutoController(mvc.Auto{
			Path:       "/role",
			Controller: new(role.Controller),
		})
		m.AutoController(mvc.Auto{
			Path:       "/admin",
			Controller: new(admin.Controller),
		})
	}
}
