package controller

import (
	"github.com/kainonly/go-bit/cipher"
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/app/system/service"
)

type Dependency struct {
	fx.In

	Crud   *crud.Crud
	Cookie *cookie.Cookie
	Cipher *cipher.Cipher

	AdminService *service.Admin
}

var Provides = fx.Provide(
	NewIndex,
	NewResource,
	NewAdmin,
)
