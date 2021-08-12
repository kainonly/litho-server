package controller

import (
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/service"
)

type Services struct {
	fx.In

	Crud   *crud.Crud
	Cooike *cookie.Cookie
	*service.Admin
}

var Provides = fx.Provide(
	NewIndex,
	NewResource,
	NewAdmin,
)
