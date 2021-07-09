package service

import (
	"context"
	"lab-api/config"
	"time"
)

type Lock struct {
	*Dependent

	key    string
	option config.Lock
}

func NewLock(dep *Dependent) *Lock {
	return &Lock{
		Dependent: dep,
		key:       dep.Config.App.Key("lock:"),
		option:    dep.Config.App.Lock,
	}
}

// Inc 增加锁定次数
func (x *Lock) Inc(ctx context.Context, key string) (err error) {
	if err = x.Redis.Incr(ctx, x.key+key).Err(); err != nil {
		return
	}
	return x.Renew(ctx, key)
}

// Renew 锁定续时
func (x *Lock) Renew(ctx context.Context, key string) error {
	return x.Redis.Expire(ctx, x.key+key, time.Second*time.Duration(x.option.RecoveryTime)).Err()
}

// Check 验证是否锁定
func (x *Lock) Check(ctx context.Context, key string) (result bool, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.key+key).Result(); err != nil {
		return
	}
	if exists == 0 {
		return true, nil
	}
	var count int64
	if count, err = x.Redis.Get(ctx, x.key+key).Int64(); err != nil {
		return
	}
	return count < x.option.Limit, nil
}

// Cancel 取消锁定
func (x *Lock) Cancel(ctx context.Context, key string) error {
	return x.Redis.Del(ctx, x.key+key).Err()
}
