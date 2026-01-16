package sessions

import (
	"context"
	"fmt"
	"server/common"
	"time"

	"github.com/goforj/wire"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	SessionsX *Service
}

type Service struct {
	*common.Inject
}

func (x *Service) Key(name string) string {
	return fmt.Sprintf(`sessions:system:%s`, name)
}

type ScanFn func(key string)

func (x *Service) Scan(ctx context.Context, fn ScanFn) {
	iter := x.RDb.Scan(ctx, 0, x.Key("*"), 0).Iterator()
	for iter.Next(ctx) {
		fn(iter.Val())
	}
}

func (x *Service) Verify(ctx context.Context, name string, jti string) bool {
	result := x.RDb.Get(ctx, x.Key(name)).Val()
	return result == jti
}

func (x *Service) Set(ctx context.Context, name string, jti string) string {
	return x.RDb.Set(ctx, x.Key(name), jti, time.Hour).Val()
}

func (x *Service) Renew(ctx context.Context, userId string) bool {
	return x.RDb.Expire(ctx, x.Key(userId), time.Hour).Val()
}
