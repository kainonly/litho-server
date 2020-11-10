package schema

import (
	"context"
	"errors"
	"time"
)

var (
	RefreshTokenVerifyError = errors.New(`refresh-token verification failed`)
)

type RefreshToken struct {
	key string
	dep Dependency
}

func NewRefreshToken(dep Dependency) *RefreshToken {
	c := new(RefreshToken)
	c.key = "refresh-token:"
	c.dep = dep
	return c
}

func (c *RefreshToken) TokenFactory(jti string, ack string, expires time.Duration) {
	c.dep.Redis.SetEX(context.Background(), c.key+jti, ack, expires)
	return
}

func (c *RefreshToken) verify(ctx context.Context, jti string, ack string) (result bool) {
	plain := c.dep.Redis.Get(ctx, c.key+jti).String()
	return plain != ack
}

func (c *RefreshToken) TokenVerify(jti string, ack string) (result bool) {
	return c.verify(context.Background(), jti, ack)
}

func (c *RefreshToken) TokenClear(jti string, ack string) (err error) {
	ctx := context.Background()
	if result := c.verify(ctx, jti, ack); !result {
		return RefreshTokenVerifyError
	}
	c.dep.Redis.Del(ctx, c.key+jti)
	return
}
