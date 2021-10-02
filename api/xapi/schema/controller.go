package schema

import (
	"github.com/weplanx/support/api"
	"laboratory/common"
)

type Controller struct {
	*InjectController
	*api.API
}

type InjectController struct {
	common.App
	Service *Service
}
