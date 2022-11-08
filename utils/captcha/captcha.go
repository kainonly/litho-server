package captcha

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type Captcha struct {
	Namespace string
	Redis     *redis.Client
}

func NewCaptcha(namespace string, client *redis.Client) *Captcha {
	return &Captcha{
		Namespace: namespace,
		Redis:     client,
	}
}

// Key 验证命名
func (x *Captcha) Key(name string) string {
	return fmt.Sprintf(`%s:captcha:%s`, x.Namespace, name)
}

// Create 创建验证码
func (x *Captcha) Create(ctx context.Context, name string, code string, ttl time.Duration) error {
	return x.Redis.
		Set(ctx, x.Key(name), code, ttl).
		Err()
}

// Exists 存在验证码
func (x *Captcha) Exists(ctx context.Context, name string) (_ bool, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.Key(name)).Result(); err != nil {
		return
	}
	return exists != 0, nil
}

var (
	ErrCaptchaNotExists    = errors.NewPublic("验证码不存在")
	ErrCaptchaInconsistent = errors.NewPublic("无效的验证码")
)

// Verify 校验验证码
func (x *Captcha) Verify(ctx context.Context, name string, code string) (err error) {
	var exists bool
	if exists, err = x.Exists(ctx, name); err != nil {
		return
	}
	if !exists {
		return ErrCaptchaNotExists
	}

	var value string
	if value, err = x.Redis.Get(ctx, x.Key(name)).Result(); err != nil {
		return
	}
	if value != code {
		return ErrCaptchaInconsistent
	}

	return
}

// Delete 移除验证码
func (x *Captcha) Delete(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.Key(name)).Err()
}
