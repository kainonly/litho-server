package schema

import (
	"context"
	"errors"
	"github.com/kainonly/gin-extra/authx"
	"time"
)

var (
	RefreshTokenVerifyError = errors.New(`refresh-token verification failed`)
)

type RefreshToken struct {
	authx.RefreshTokenAPI

	key string
	Dependency
}

func NewRefreshToken(dep Dependency) *RefreshToken {
	c := new(RefreshToken)
	c.key = "refresh-token:"
	c.Dependency = dep
	return c
}

func (c *RefreshToken) Factory(value ...interface{}) {
	c.Redis.SetEX(
		context.Background(),
		c.key+value[0].(string),
		value[1].(string),
		value[2].(time.Duration),
	)
	return
}

func (c *RefreshToken) Renewal(value ...interface{}) {
	c.Redis.Expire(context.Background(), value[0].(string), value[1].(time.Duration))
}

func (c *RefreshToken) Verify(value ...interface{}) (result bool) {
	return c.tokenVerify(context.Background(), value[0].(string), value[1].(string))
}

func (c *RefreshToken) tokenVerify(ctx context.Context, jti string, ack string) (result bool) {
	plain := c.Redis.Get(ctx, c.key+jti).String()
	return plain != ack
}

func (c *RefreshToken) Destory(value ...interface{}) (err error) {
	ctx := context.Background()
	if result := c.tokenVerify(ctx, value[0].(string), value[1].(string)); !result {
		return RefreshTokenVerifyError
	}
	c.Redis.Del(ctx, c.key+value[0].(string))
	return
}
