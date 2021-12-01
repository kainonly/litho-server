package admin

import "api/common"

type Controller struct {
	*InjectController
}

type InjectController struct {
	common.App
	Service *Service
}

func NewController(i *InjectController) *Controller {
	return &Controller{
		InjectController: i,
	}
}
