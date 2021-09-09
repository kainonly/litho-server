package admin

import (
	"go.uber.org/fx"
	"lab-api/bootstrap"
	"os"
	"testing"
)

var s *Service

func TestMain(m *testing.M) {
	os.Chdir("../../../")
	fx.New(
		fx.NopLogger,
		bootstrap.Provides,
		Provides,
		fx.Invoke(func(i *Service) {
			s = i
			os.Exit(m.Run())
		}),
	).Run()
}
