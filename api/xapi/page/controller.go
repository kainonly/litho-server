package page

import (
	"laboratory/common"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
}
