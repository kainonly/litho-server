package controller

import (
	"github.com/kainonly/go-bit/cookie"
	"github.com/kainonly/go-bit/crud"
	"go.uber.org/fx"
	"lab-api/app/system/service"
)

type Dependency struct {
	fx.In

	*crud.Crud
	*cookie.Cookie

	AdminService *service.Admin
}

var Provides = fx.Provide(
	NewResource,
	NewAdmin,
)
