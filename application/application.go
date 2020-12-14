package application

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/helper/token"
	"github.com/kainonly/gin-extra/mvcx"
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

func Application(router *gin.Engine, dependency common.Dependency) {
	cfg := dependency.Config
	token.Key = []byte(cfg.App.Key)
	router.GET("/", routes.Default)
	system := router.Group("/system")
	{
		mvc := mvcx.Initialize(system, dependency)
		mvc.AutoController("/main", &controller.Controller{}, mvcx.Middleware{
			Handle: auth.Auth(),
			Only:   []string{"Resource"},
		})
		mvc.AutoController("/acl", &acl.Controller{})
		mvc.AutoController("/resource", &resource.Controller{})
		mvc.AutoController("/policy", &policy.Controller{})
		mvc.AutoController("/role", &role.Controller{})
		mvc.AutoController("/admin", &admin.Controller{})
	}
}
