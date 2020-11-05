package application

import (
	"github.com/kainonly/gin-helper/mvc"
	"taste-api/application/common"
	"taste-api/application/controller/admin"
	"taste-api/routes"
)

func Application(mvc *mvc.Mvc, dependency common.Dependency) {
	mvc.Dependency(&dependency)
	mvc.GET("/", mvc.Handle(routes.Default))
	mvc.AutoController("/admin", new(admin.Controller))
}
