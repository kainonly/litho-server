package sessions

import (
	"api/app/vars"
	"api/common"
	"context"
	"strings"
)

type Service struct {
	*common.Inject
	Vars *vars.Service
}

// Gets 获取所有会话
func (x *Service) Gets(ctx context.Context) (values []string, err error) {
	var cursor uint64
	for {
		var keys []string
		var next uint64
		if keys, next, err = x.Redis.Scan(ctx,
			cursor, x.Values.KeyName("sessions", "*"), 1000,
		).Result(); err != nil {
			return
		}
		uids := make([]string, len(keys))
		for k, v := range keys {
			uids[k] = strings.Replace(v, x.Values.KeyName("sessions", ""), "", -1)
		}
		values = append(values, uids...)
		if next == 0 {
			break
		}
		cursor = next
	}
	return
}

// Verify 验证会话一致性
func (x *Service) Verify(ctx context.Context, uid string, jti string) (_ bool, err error) {
	var value string
	if value, err = x.Redis.Get(ctx, x.Values.KeyName("sessions", uid)).Result(); err != nil {
		return
	}
	return value == jti, nil
}

// Set 设置会话
func (x *Service) Set(ctx context.Context, uid string, jti string) (err error) {
	expiration := x.Vars.GetExpiration(ctx)
	if err = x.Redis.Set(ctx, x.Values.KeyName("sessions", uid), jti, expiration).Err(); err != nil {
		return
	}
	return
}

// Renew 续约会话
func (x *Service) Renew(ctx context.Context, uid string) (err error) {
	expiration := x.Vars.GetExpiration(ctx)
	if err = x.Redis.Expire(ctx,
		x.Values.KeyName("sessions", uid), expiration,
	).Err(); err != nil {
		return
	}
	return
}

// Delete 删除会话
func (x *Service) Delete(ctx context.Context, uid string) (err error) {
	return x.Redis.Del(ctx, x.Values.KeyName("sessions", uid)).Err()
}

// BulkDelete 删除所有会话
func (x *Service) BulkDelete(ctx context.Context) (err error) {
	var cursor uint64
	var keys []string
	for {
		var next uint64
		if keys, next, err = x.Redis.Scan(ctx,
			cursor, x.Values.KeyName("sessions", "*"), 1000,
		).Result(); err != nil {
			return
		}
		if next == 0 {
			break
		}
		cursor = next
	}
	return x.Redis.Del(ctx, keys...).Err()
}
