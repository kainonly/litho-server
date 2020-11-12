package common

import (
	curd "github.com/kainonly/gin-curd"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"taste-api/application/cache"
	"taste-api/config"
)

type Dependency struct {
	fx.In

	Config *config.Config
	Db     *gorm.DB
	Cache  *cache.Cache
	Curd   *curd.Curd
}
