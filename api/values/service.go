package values

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/server/common"
	"time"
)

type Service struct {
	common.Inject
}

// Key 命名
func (x *Service) Key() string {
	return x.Values.Name("values")
}

// Load 载入配置
func (x *Service) Load(ctx context.Context) (err error) {
	var b []byte
	b, err = x.Store.GetBytes("values")
	if err != nil {
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
			if _, err = x.Store.PutBytes("values", b); err != nil {
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
		if o == nil || o.ModTime.Unix() < current.Unix() {
			continue
		}
	}

	return
}

// Get 获取动态配置
func (x *Service) Get(keys ...string) (data map[string]interface{}, err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values"); err != nil {
		return
	}
	if err = sonic.Unmarshal(b, &data); err != nil {
		return
	}
	fmt.Println(data)
	sets := make(map[string]bool)
	for _, key := range keys {
		sets[key] = true
	}
	all := len(sets) == 0
	data = make(map[string]interface{})
	for k, v := range data {
		if !all && !sets[k] {
			continue
		}
		if secret[k] {
			// 密文
			if v != nil || v != "" {
				data[k] = "*"
			} else {
				data[k] = "-"
			}
		} else {
			data[k] = v
		}
	}
	return
}

// Set 设置动态配置
func (x *Service) Set(ctx context.Context, data map[string]interface{}) (err error) {
	//
	//// 合并覆盖
	//for k, v := range data {
	//	x.Values.DynamicValues[k] = v
	//}
	//return x.Update(ctx)
	return
}

// Remove 移除动态配置
func (x *Service) Remove(ctx context.Context, key string) (err error) {
	//delete(x.Values.DynamicValues, key)
	//return x.Update(ctx)
	return
}

// Update 更新配置
func (x *Service) Update(ctx context.Context) (err error) {
	var b []byte
	if b, err = sonic.Marshal(x.Values.DynamicValues); err != nil {
		return
	}
	if err = x.Redis.Set(ctx, x.Key(), b, 0).Err(); err != nil {
		return
	}

	// 发布同步配置
	if err = x.Nats.Publish(x.Key(), []byte("sync")); err != nil {
		return
	}

	return
}
