package controller

import (
	"github.com/kainonly/go-bit/cookie"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"lab-api/app/api/service"
	"lab-api/common"
)

type Dependency struct {
	fx.In

	App    *common.App
	Db     *gorm.DB
	Cookie *cookie.Cookie

	IndexService *service.Index
}

var Provides = fx.Provide(
	NewIndex,
	NewDeveloper,
)
