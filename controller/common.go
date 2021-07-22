package controller

import (
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	NewIndex,
)
