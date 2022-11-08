package locker

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Locker struct {
	Namespace string
	Redis     *redis.Client
}

func NewLocker(namespace string, client *redis.Client) *Locker {
	return &Locker{
		Namespace: namespace,
		Redis:     client,
	}
}

// Key 锁定命名
func (x *Locker) Key(name string) string {
	return fmt.Sprintf(`%s:locker:%s`, x.Namespace, name)
}

// Update 更新锁定
func (x *Locker) Update(ctx context.Context, name string, ttl time.Duration) (err error) {
	key := x.Key(name)
	var exists int64
	if exists, err = x.Redis.
		Exists(ctx, key).
		Result(); err != nil {
		return
	}

	if exists == 0 {
		if err = x.Redis.
			Set(ctx, key, 1, ttl).
			Err(); err != nil {
			return
		}
	} else {
		if err = x.Redis.
			Incr(ctx, key).
			Err(); err != nil {
			return
		}
	}
	return
}

// Verify 校验锁定，True 为锁定
func (x *Locker) Verify(ctx context.Context, name string, n int64) (result bool, err error) {
	key := x.Key(name)
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if exists == 0 {
		return
	}

	var count int64
	if count, err = x.Redis.Get(ctx, key).Int64(); err != nil {
		return
	}

	return count >= n, nil
}

// Delete 移除锁定
func (x *Locker) Delete(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.Key(name)).Err()
}
