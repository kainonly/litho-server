package application

import (
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

func Application(mvc *mvc.Mvc, dependency common.Dependency) {
	mvc.Dependency(&dependency)
	mvc.GET("/", mvc.Handle(routes.Default))
	mvc.AutoController("/main", new(controller.Controller))
	mvc.AutoController("/acl", new(acl.Controller))
	mvc.AutoController("/resource", new(resource.Controller))
	mvc.AutoController("/policy", new(policy.Controller))
	mvc.AutoController("/role", new(role.Controller))
	mvc.AutoController("/admin", new(admin.Controller))
}
