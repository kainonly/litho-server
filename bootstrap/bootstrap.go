package bootstrap

import (
	"context"
	"database/sql"
	"github.com/caarlos0/env/v8"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/hertz-contrib/requestid"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/weplanx/server/common"
	"go.uber.org/fx"
	"os"
)

var Provides = fx.Provide(
	LoadStaticValues,
	UseDatabase,
	UseRedis,
	UseHertz,
)

func LoadStaticValues() (v *common.Values, err error) {
	v = new(common.Values)
	if err = env.Parse(v); err != nil {
		return
	}
	return
}

func UseDatabase(v *common.Values) (db *bun.DB) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(v.Database.Url)))
	db = bun.NewDB(sqldb, pgdialect.New())
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
