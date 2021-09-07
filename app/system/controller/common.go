package controller

import (
	"github.com/kainonly/go-bit/cipher"
	"github.com/kainonly/go-bit/cookie"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"lab-api/app/system/service"
)

type Dependency struct {
	fx.In

	Db     *gorm.DB
	Cookie *cookie.Cookie
	Cipher *cipher.Cipher

	IndexService    *service.Index
	ResourceService *service.Resource
}

var Provides = fx.Provide(
	NewIndex,
	NewResource,
)
