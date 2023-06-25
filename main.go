package main

import (
	"github.com/weplanx/server/admin"
	"github.com/weplanx/server/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		bootstrap.Provides,
		admin.Options,
	).Run()
}
