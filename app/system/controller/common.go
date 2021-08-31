package controller

import (
	"github.com/google/wire"
	"github.com/kainonly/go-bit/cookie"
	"gorm.io/gorm"
	"lab-api/app/system/service"
)

type Dependency struct {
	Db     *gorm.DB
	Cookie *cookie.Cookie

	IndexService *service.Index
}

var Provides = wire.NewSet(
	wire.Struct(new(Dependency), "*"),
	NewIndex,
	NewAdmin,
)
