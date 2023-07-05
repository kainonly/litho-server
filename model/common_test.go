package model_test

import (
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"os"
	"testing"
)

type Inject struct {
	fx.In

	V   *common.Values
	Mgo *mongo.Client
	Db  *mongo.Database
}

var x *Inject

func TestMain(m *testing.M) {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.LoadStaticValues,
			bootstrap.UseMongoDB,
			bootstrap.UseRedis,
		),
		fx.Invoke(func(i Inject) {
			x = &i
			os.Exit(m.Run())
		}),
	).Run()
}
