package roles

import (
	"api/common"
	"github.com/weplanx/go/api"
)

type Controller struct {
	*InjectController
	*api.Controller
}

type InjectController struct {
	common.Inject
	APIs    *api.API
	Service *Service
}
