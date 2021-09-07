package controller

import (
	"github.com/kainonly/go-bit/cookie"
	"go.uber.org/fx"
	"lab-api/app/api/service"
)

type Dependency struct {
	fx.In

	*cookie.Cookie

	IndexService *service.Index
}

var Provides = fx.Provide(
	NewIndex,
)
