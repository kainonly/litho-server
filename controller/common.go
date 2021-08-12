package controller

import (
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/service"
)

type Dependency struct {
	fx.In

	*crud.Crud
	*cookie.Cookie
	*service.Admin
}

var Provides = fx.Provide(
	NewIndex,
	NewResource,
	NewAdmin,
)
