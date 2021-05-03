package schema

import (
	"context"
	"time"
)

type Lock struct {
	key   string
	limit int64
	renew time.Duration
	Dependency
}

func NewLock(dep Dependency, limit int64, renew time.Duration) *Lock {
	c := new(Lock)
	c.key = "system:lock:"
	c.Dependency = dep
	c.limit = limit
	c.renew = renew
	return c
}

func (c *Lock) Remove(ctx context.Context, name string) {
	c.Redis.Del(ctx, c.key+name)
}

func (c *Lock) Check(ctx context.Context, name string) bool {
	if exists := c.Redis.Exists(ctx, c.key+name).Val(); exists == 0 {
		return true
	}
	count, _ := c.Redis.Get(ctx, c.key+name).Int64()
	return count < c.limit
}

func (c *Lock) Inc(ctx context.Context, name string) {
	c.Redis.Incr(ctx, c.key+name)
	c.Lock(ctx, name)
}

func (c *Lock) Lock(ctx context.Context, name string) {
	c.Redis.Expire(ctx, c.key+name, c.renew)
}
