package page

import (
	"api/common"
)

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.Inject
	Service *Service
}
