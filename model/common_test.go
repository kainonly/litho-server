package model

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
	"lab-api/bootstrap"
	"os"
	"testing"
)

var tx *gorm.DB

func TestMain(m *testing.M) {
	os.Chdir("../")
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.LoadSettings,
			bootstrap.InitializeDatabase,
		),
		fx.Invoke(func(db *gorm.DB) {
			tx = db
			os.Exit(m.Run())
		}),
	).Run()
}
