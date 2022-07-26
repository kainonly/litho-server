package captcha

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
	return x.Values.Key("captcha", name)
}

// Create 创建验证码
func (x *Service) Create(ctx context.Context, name string, code string, ttl time.Duration) error {
	return x.Redis.
		Set(ctx, x.Key(name), code, ttl).
		Err()
}

// Exists 存在验证码
func (x *Service) Exists(ctx context.Context, name string) (_ bool, err error) {
	var exists int64
	if exists, err = x.Redis.
		Exists(ctx, x.Key(name)).
		Result(); err != nil {

	}
	return exists != 0, nil
}

// Verify 校验验证码
func (x *Service) Verify(ctx context.Context, name string, code string) (_ bool, err error) {
	var value string
	if value, err = x.Redis.
		Get(ctx, x.Key(name)).
		Result(); err != nil {
		return
	}
	return value == code, nil
}

// Delete 移除验证码
func (x *Service) Delete(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.Key(name)).Err()
}
