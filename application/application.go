package application

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kainonly/gin-extra/authx"
	"github.com/kainonly/gin-extra/mvcx"
	"github.com/kainonly/gin-extra/tokenx"
	"github.com/kainonly/gin-extra/typ"
	"lab-api/application/common"
	"lab-api/application/controller"
	"lab-api/application/controller/acl"
	"lab-api/application/controller/admin"
	"lab-api/application/controller/permission"
	"lab-api/application/controller/policy"
	"lab-api/application/controller/resource"
	"lab-api/application/controller/role"
	"lab-api/routes"
	"net/http"
	"strings"
)

func Application(router *gin.Engine, dependency common.Dependency) {
	cfg := dependency.Config
	tokenx.LoadKey([]byte(cfg.App.Key))
	router.GET("/", routes.Default)
	system := router.Group("/system")
	{
		auth := authx.Middleware(typ.Cookie{
			Name:     "system",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}, dependency.Redis.RefreshToken)
		unifyMiddleware := []mvcx.Middleware{
			{
				Handle: auth,
			},
			{
				Handle: func(ctx *gin.Context) {
					path := strings.Replace(ctx.Request.URL.Path, "/system/", "", 1)
					acts := strings.Split(path, "/")
					var err error
					var auth jwt.MapClaims
					if auth, err = authx.Get(ctx); err != nil {
						ctx.AbortWithStatusJSON(200, gin.H{
							"error": 1,
							"msg":   err.Error(),
						})
					}
					userData := dependency.Redis.Admin.Get(auth["user"].(string))
					aclSet := dependency.Redis.Role.Get(
						strings.Split(userData["role"].(string), ","),
						"acl",
					)
					if userData["acl"] != nil {
						aclSet.Add(userData["acl"].([]interface{})...)
					}
					policyCursor := ""
					policyValues := []string{"0", "1"}
					for _, val := range policyValues {
						if aclSet.Contains(acts[0] + ":" + val) {
							policyCursor = val
						}
					}
					if policyCursor == "" {
						ctx.AbortWithStatusJSON(200, gin.H{
							"error": 1,
							"msg":   "rbac invalid, policy is empty",
						})
					}
					scope := dependency.Redis.Acl.Get(acts[0], policyCursor)
					if scope.Empty() {
						ctx.AbortWithStatusJSON(200, gin.H{
							"error": 1,
							"msg":   "rbac invalid, scope is empty",
						})
					}
					if !scope.Contains(acts[1]) {
						ctx.AbortWithStatusJSON(200, gin.H{
							"error": 1,
							"msg":   "rbac invalid, access denied",
						})
					}
					ctx.Next()
				},
			},
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
	}
}
