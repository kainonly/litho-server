package planx

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/kainonly/go-bit/support"
	"lab-api/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}

func (x *Service) planxKey() string {
	return x.Set.RedisKey("planx")
}

func (x *Service) GetFromCache(ctx context.Context) (data []support.Planx, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.planxKey()).Result(); err != nil {
		return
	}
	if exists == 0 {
		if err = x.RefreshCache(ctx); err != nil {
			return
		}
	}
	var value string
	if value, err = x.Redis.Get(ctx, x.planxKey()).Result(); err != nil {
		return
	}
	if err = jsoniter.Unmarshal([]byte(value), &data); err != nil {
		return
	}
	return
}

func (x *Service) RefreshCache(ctx context.Context) (err error) {
	var data []support.Planx
	if err = x.Db.WithContext(ctx).
		Model(&support.Planx{}).
		Find(&data).Error; err != nil {
		return
	}
	var value []byte
	if value, err = jsoniter.Marshal(&data); err != nil {
		return
	}
	if err = x.Redis.Set(ctx, x.planxKey(), value, 0).Err(); err != nil {
		return
	}
	return
}

func (x *Service) RemoveCache(ctx context.Context) error {
	return x.Redis.Del(ctx, x.planxKey()).Err()
}
