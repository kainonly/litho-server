package mock

import (
	"api/bootstrap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"os"
	"testing"
)

var (
	Db *mongo.Database
)

func TestMain(m *testing.M) {
	os.Chdir("../")
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.SetValues,
			bootstrap.UseDatabase,
		),
		fx.Invoke(func(client *mongo.Client, db *mongo.Database) {
			Db = db
			os.Exit(m.Run())
		}),
	)
}
