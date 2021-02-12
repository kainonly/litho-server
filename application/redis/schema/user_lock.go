package schema

import (
	"context"
	"time"
)

type UserLock struct {
	key string
	Dependency
}

func NewUserLock(dep Dependency) *UserLock {
	c := new(UserLock)
	c.key = "system:user-lock:"
	c.Dependency = dep
	return c
}

func (c *UserLock) Remove(username string) {
	c.Redis.Del(context.Background(), c.key+username)
}

func (c *UserLock) Check(username string) bool {
	ctx := context.Background()
	exists := c.Redis.Exists(ctx, c.key+username).Val()
	if exists == 0 {
		return true
	}
	count, err := c.Redis.Get(ctx, c.key+username).Int64()
	if err != nil {
		return false
	}
	return count < 5
}

func (c *UserLock) Inc(ctx context.Context, username string) {
	c.Redis.Incr(ctx, c.key+username)
	c.Lock(ctx, username)
}

func (c *UserLock) Lock(ctx context.Context, username string) {
	c.Redis.Expire(ctx, c.key+username, time.Second*900)
}
