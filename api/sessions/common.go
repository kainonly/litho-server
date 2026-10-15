package sessions

import (
	"context"
	"fmt"
	"litho-api/common"
	"time"

	"github.com/goforj/wire"
	"github.com/kainonly/go/help"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

const Key = "sessions"

type Controller struct {
	V *common.Values

	SessionsX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]any

var ErrAccess = help.E(0, `您不具备系统权限，请联系超级管理员进行管理`)

func (x *Service) CheckAccess(ctx context.Context, user *common.IAMUser, ids ...string) (err error) {
	if err = user.Can(`WM`); err != nil {
		return ErrAccess
	}
	return
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
