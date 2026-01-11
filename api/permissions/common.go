package permissions

import (
	"server/common"

	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{PermissionsX: i} },
	func(i common.Inject) *Service { return &Service{Inject: &i} },
)

type Controller struct {
	PermissionsX *Service
}

type Service struct {
	*common.Inject
}
