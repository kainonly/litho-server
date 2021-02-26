package application

import (
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/mvcx"
	"github.com/kainonly/gin-extra/rbacx"
	"github.com/kainonly/gin-extra/storage/cos"
	"github.com/kainonly/gin-extra/tokenx"
	"lab-api/application/common"
	"lab-api/application/controller"
	"lab-api/application/controller/acl"
	"lab-api/application/controller/admin"
	"lab-api/application/controller/permission"
	"lab-api/application/controller/picture"
	"lab-api/application/controller/picture_type"
	"lab-api/application/controller/policy"
	"lab-api/application/controller/resource"
	"lab-api/application/controller/role"
	"lab-api/routes"
)

func Application(router *gin.Engine, dependency common.Dependency) {
	cfg := dependency.Config
	tokenx.LoadKey([]byte(cfg.App.Key))
	switch cfg.Storage.Type {
	case "cos":
		cos.DEF = cos.Option{
			Bucket:    cfg.Storage.Option["bucket"].(string),
			Region:    cfg.Storage.Option["region"].(string),
			SecretID:  cfg.Storage.Option["secret_id"].(string),
			SecretKey: cfg.Storage.Option["secret_key"].(string),
		}
		break
	}
	router.GET("/", routes.Default)
	system := router.Group("/system")
	{
		auth := authx.Middleware(common.SystemCookie, dependency.Redis.RefreshToken)
		rbac := rbacx.Middleware("/system/",
			dependency.Redis.Admin,
			dependency.Redis.Role,
			dependency.Redis.Acl,
		)
		unifyMiddleware := []mvcx.Middleware{
			{Handle: auth}, {Handle: rbac},
		}
		mvc := mvcx.Initialize(system, dependency)
		mvc.AutoController("/main", &controller.Controller{}, mvcx.Middleware{
			Handle: auth,
			Only:   []string{"Resource", "Information", "Update"},
		})
		mvc.AutoController("/acl", &acl.Controller{}, unifyMiddleware...)
		mvc.AutoController("/resource", &resource.Controller{}, unifyMiddleware...)
		mvc.AutoController("/policy", &policy.Controller{}, unifyMiddleware...)
		mvc.AutoController("/permission", &permission.Controller{}, unifyMiddleware...)
		mvc.AutoController("/role", &role.Controller{}, unifyMiddleware...)
		mvc.AutoController("/admin", &admin.Controller{}, unifyMiddleware...)
		mvc.AutoController("/picture_type", &picture_type.Controller{}, unifyMiddleware...)
		mvc.AutoController("/picture", &picture.Controller{}, unifyMiddleware...)
	}
}
