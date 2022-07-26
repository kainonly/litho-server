package locker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/server/common"
	"time"
)

type Service struct {
	Values *common.Values
	Redis  *redis.Client
}

func (x *Service) Key(name string) string {
	return x.Values.Key("locker", name)
}

// Update 更新锁定
func (x *Service) Update(ctx context.Context, name string, ttl time.Duration) (err error) {
	var exists int64
	if exists, err = x.Redis.
		Exists(ctx, x.Key(name)).
		Result(); err != nil {
		return
	}

	if exists == 0 {
		if err = x.Redis.
			Set(ctx, x.Key(name), 1, ttl).
			Err(); err != nil {
			return
		}
	} else {
		if err = x.Redis.
			Incr(ctx, x.Key(name)).
			Err(); err != nil {
			return
		}
	}
	return
}

// Verify 校验锁定，True 为锁定
func (x *Service) Verify(ctx context.Context, name string, n int64) (result bool, err error) {
	var count int64
	if count, err = x.Redis.
		Get(ctx, x.Key(name)).
		Int64(); err != nil {
		return
	}

	return count >= n, nil
}

// Delete 移除锁定
func (x *Service) Delete(ctx context.Context, name string) {
	x.Redis.Del(ctx, x.Key(name))
}
