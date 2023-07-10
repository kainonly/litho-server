package bootstrap

import (
	"context"
	"github.com/caarlos0/env/v9"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/requestid"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.uber.org/fx"
	"os"
)

func LoadStaticValues() (v *common.Values, err error) {
	v = new(common.Values)
	if err = env.Parse(v); err != nil {
		return
	}
	return
}

func UseMongoDB(v *common.Values) (client *mongo.Client, db *mongo.Database, err error) {
	if client, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(v.Database.Url),
	); err != nil {
		return
	}
	option := options.Database().
		SetWriteConcern(writeconcern.Majority())
	db = client.Database(v.Database.Name, option)
	return
}

// UseRedis
// https://github.com/go-redis/redis
func UseRedis(v *common.Values) (client *redis.Client, err error) {
	opts, err := redis.ParseURL(v.Database.Redis)
	if err != nil {
		return
	}
	client = redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return
	}
	return
}

// UseHertz
// https://www.cloudwego.io/zh/docs/hertz/reference/config
func UseHertz(lc fx.Lifecycle, v *common.Values) (h *server.Hertz, err error) {
	opts := []config.Option{
		server.WithHostPorts(v.Address),
	}

	if os.Getenv("MODE") != "release" {
		opts = append(opts, server.WithExitWaitTime(0))
	}

	opts = append(opts)

	h = server.Default(opts...)

	h.Use(
		requestid.New(),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go h.Spin()
			return nil
		},
	})

	return
}
