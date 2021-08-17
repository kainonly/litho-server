package service

import (
	"go.uber.org/fx"
	"lab-api/bootstrap"
	"os"
	"testing"
)

var s *Tests

type Tests struct {
	fx.In

	Resource *Resource
}

func TestMain(m *testing.M) {
	os.Chdir(`../../../`)
	fx.New(
		fx.Provide(
			bootstrap.LoadConfiguration,
			bootstrap.InitializeDatabase,
			bootstrap.InitializeRedis,
		),
		Provides,
		fx.Invoke(func(tests Tests) {
			s = &tests
			os.Exit(m.Run())
		}),
	).Run()
}
