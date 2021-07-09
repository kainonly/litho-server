package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"lab-api/config"
	"time"
)

type Lock struct {
	key    string
	option *config.Lock
	redis  *redis.Client
}

func NewLock(cfg *config.Config, redis *redis.Client) *Lock {
	return &Lock{
		key:    cfg.App.Key("lock:"),
		option: &cfg.App.Lock,
		redis:  redis,
	}
}

// Inc 增加锁定次数
func (x *Lock) Inc(ctx context.Context, key string) (err error) {
	if err = x.redis.Incr(ctx, x.key+key).Err(); err != nil {
		return
	}
	return x.Renew(ctx, key)
}

// Renew 锁定续时
func (x *Lock) Renew(ctx context.Context, key string) error {
	return x.redis.Expire(ctx, x.key+key, time.Second*time.Duration(x.option.RecoveryTime)).Err()
}

// Check 验证是否锁定
func (x *Lock) Check(ctx context.Context, key string) (result bool, err error) {
	var exists int64
	if exists, err = x.redis.Exists(ctx, x.key+key).Result(); exists == 0 {
		return
	}
	var count int64
	if count, err = x.redis.Get(ctx, x.key+key).Int64(); err != nil {
		return
	}
	return count < x.option.RecoveryTime, nil
}

// Cancel 取消锁定
func (x *Lock) Cancel(ctx context.Context, key string) error {
	return x.redis.Del(ctx, x.key+key).Err()
}
