package system

import (
	"context"
	"laboratory/common"
	"time"
)

type InjectService struct {
	*common.App
}

type Service struct {
	*InjectService
}

func NewService(i *InjectService) *Service {
	return &Service{
		InjectService: i,
	}
}

func (x *Service) verifyCodeKey(name string) string {
	return x.Set.RedisKey("verify:" + name)
}

func (x *Service) CreateVerifyCode(ctx context.Context, name string, code string) error {
	return x.Redis.Set(ctx, x.verifyCodeKey(name), code, time.Minute).Err()
}

// VerifyCode 校验验证码
func (x *Service) VerifyCode(ctx context.Context, name string, code string) (result bool, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.verifyCodeKey(name)).Result(); err != nil {
		return
	}
	if exists == 0 {
		return false, nil
	}
	var value string
	if value, err = x.Redis.Get(ctx, x.verifyCodeKey(name)).Result(); err != nil {
		return
	}
	return value == code, nil
}

// RemoveVerifyCode 移除验证码
func (x *Service) RemoveVerifyCode(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.verifyCodeKey(name)).Err()
}
