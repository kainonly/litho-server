package service

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Query func(tx *gorm.DB) *gorm.DB

var Provides = fx.Provide(
	NewAdmins,
)
