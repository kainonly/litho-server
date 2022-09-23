package values

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/server/common"
	"time"
)

type Service struct {
	common.Inject
}

// Load 载入配置
func (x *Service) Load(ctx context.Context) (err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		// 不存在配置则初始化
		if errors.Is(err, nats.ErrObjectNotFound) {
			x.Values.DynamicValues = common.DynamicValues{
				SessionTTL:      time.Hour,
				LoginTTL:        time.Minute * 15,
				LoginFailures:   5,
				IpLoginFailures: 10,
				IpWhitelist:     []string{},
				IpBlacklist:     []string{},
				PwdStrategy:     1,
				PwdTTL:          time.Hour * 24 * 365,
			}
			if b, err = sonic.Marshal(x.Values.DynamicValues); err != nil {
				return
			}
			if _, err = x.Store.PutBytes("values", b, nats.Context(ctx)); err != nil {
				return
			}
			return nil
		}
		return
	}

	if err = sonic.Unmarshal(b, &x.Values.DynamicValues); err != nil {
		return
	}

	return
}

// Sync 同步节点动态配置
func (x *Service) Sync(ctx context.Context) (err error) {
	if err = x.Load(ctx); err != nil {
		return
	}

	var watch nats.ObjectWatcher
	if watch, err = x.Store.Watch(); err != nil {
		return
	}
	current := time.Now()
	for o := range watch.Updates() {
		if o == nil || o.Name != "values" || o.ModTime.Unix() < current.Unix() {
			continue
		}
		// 同步动态配置
		var b []byte
		if b, err = x.Store.GetBytes("values"); err != nil {
			return
		}
		if err = sonic.Unmarshal(b, &x.Values.DynamicValues); err != nil {
			return
		}
	}

	return
}

// Get 获取动态配置
func (x *Service) Get(ctx context.Context, keys ...string) (values map[string]interface{}, err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		return
	}
	if err = sonic.Unmarshal(b, &values); err != nil {
		return
	}
	sets := make(map[string]bool)
	for _, key := range keys {
		sets[key] = true
	}
	all := len(sets) == 0
	for k, v := range values {
		if !all && !sets[k] {
			continue
		}
		if secret[k] {
			// 密文
			if v != nil || v != "" {
				values[k] = "*"
			} else {
				values[k] = "-"
			}
		} else {
			values[k] = v
		}
	}
	return
}

// Set 设置动态配置
func (x *Service) Set(ctx context.Context, update map[string]interface{}) (err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		return
	}
	var values map[string]interface{}
	if err = sonic.Unmarshal(b, &values); err != nil {
		return
	}
	for k, v := range update {
		values[k] = v
	}
	return x.Update(ctx, values)
}

// Remove 移除动态配置
func (x *Service) Remove(ctx context.Context, key string) (err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		return
	}
	var values map[string]interface{}
	if err = sonic.Unmarshal(b, &values); err != nil {
		return
	}
	delete(values, key)
	return x.Update(ctx, values)
}

// Update 更新配置
func (x *Service) Update(ctx context.Context, values map[string]interface{}) (err error) {
	var b []byte
	if b, err = sonic.Marshal(values); err != nil {
		return
	}
	if _, err = x.Store.PutBytes("values", b, nats.Context(ctx)); err != nil {
		return
	}
	return
}
