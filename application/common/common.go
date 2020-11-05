package common

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"taste-api/config"
)

type Dependency struct {
	fx.In

	Config *config.Config
	Db     *gorm.DB
}

func Inject(i interface{}) *Dependency {
	return i.(*Dependency)
}
