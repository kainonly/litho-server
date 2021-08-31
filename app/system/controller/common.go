package controller

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit/cipher"
	"github.com/kainonly/go-bit/cookie"
	"gorm.io/gorm"
	"lab-api/app/system/service"
)

type Dependency struct {
	Db     *gorm.DB
	Cookie *cookie.Cookie
	Cipher *cipher.Cipher

	IndexService    *service.Index
	ResourceService *service.Resource
	AdminService    *service.Admin
}

var Provides = wire.NewSet(
	wire.Struct(new(Dependency), "*"),
	NewIndex,
	NewResource,
	NewRole,
	NewAdmin,
)
