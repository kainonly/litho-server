package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"lab-api/app/system/index"
	"lab-api/app/system/resource"
	"lab-api/common"
)

var Options = fx.Options(
	index.Provides,
	resource.Provides,
	fx.Invoke(Routes),
)

type Inject struct {
	common.App
	Index    *index.Controller
	Resource *resource.Controller
}

func Routes(r *gin.Engine, i Inject) {
}
