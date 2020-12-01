package application

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/token"
	"github.com/kainonly/gin-extra/mvc"
	"lab-api/application/common"
	"lab-api/application/controller"
	"lab-api/application/controller/acl"
	"lab-api/application/controller/admin"
	"lab-api/application/controller/policy"
	"lab-api/application/controller/resource"
	"lab-api/application/controller/role"
	"lab-api/application/middleware/auth"
	"lab-api/routes"
)

func Application(router *gin.Engine, dep common.Dependency) {
	cfg := dep.Config
	token.Key = []byte(cfg.App.Key)
	router.GET("/", routes.Default)
	system := router.Group("/system")
	{
		m := mvc.Factory(system)
		m.AutoController(mvc.Auto{
			Path:       "/main",
			Controller: &controller.Controller{Dependency: &dep},
			Middlewares: []mvc.Middleware{
				{
					Handle: auth.Auth(),
					Only:   []string{"Resource"},
				},
			},
		})
		m.AutoController(mvc.Auto{
			Path:       "/acl",
			Controller: &acl.Controller{Dependency: &dep},
		})
		m.AutoController(mvc.Auto{
			Path:       "/resource",
			Controller: &resource.Controller{Dependency: &dep},
		})
		m.AutoController(mvc.Auto{
			Path:       "/policy",
			Controller: &policy.Controller{Dependency: &dep},
		})
		m.AutoController(mvc.Auto{
			Path:       "/role",
			Controller: &role.Controller{Dependency: &dep},
		})
		m.AutoController(mvc.Auto{
			Path:       "/admin",
			Controller: &admin.Controller{Dependency: &dep},
		})
	}
}
