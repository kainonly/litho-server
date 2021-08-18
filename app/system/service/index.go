package service

import (
	"context"
	"time"
)

type Index struct {
	*Dependency
	Key string
}

func NewIndex(d Dependency) *Index {
	return &Index{
		Dependency: &d,
		Key:        d.Config.RedisKey("code:"),
	}
}

// GenerateCode 生成验证码
func (x *Index) GenerateCode(ctx context.Context, index string, code string) error {
	return x.Redis.Set(ctx, x.Key+index, code, time.Minute).Err()
}

// VerifyCode 校验验证码
func (x *Index) VerifyCode(ctx context.Context, index string, code string) (result bool, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.Key+index).Result(); err != nil {
		return
	}
	if exists == 0 {
		return false, nil
	}
	var value string
	if value, err = x.Redis.Get(ctx, x.Key+index).Result(); err != nil {
		return
	}
	return value != code, nil
}

// RemoveCode 移除验证码
func (x *Index) RemoveCode(ctx context.Context, index string) error {
	return x.Redis.Del(ctx, x.Key+index).Err()
}
