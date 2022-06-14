package values

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/thoas/go-funk"
)

type Service struct {
	Store nats.ObjectStore
}

// Get 获取配置
func (x *Service) Get(ctx context.Context, keys []string) (data map[string]interface{}, err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		return
	}
	values := make(map[string]interface{})
	if err = jsoniter.Unmarshal(b, &values); err != nil {
		return
	}
	data = make(map[string]interface{})
	for k, v := range values {
		if x.IsSecret(k) {
			// 存在数值
			if v != nil || v != "" {
				data[k] = "*"
			}
		} else {
			data[k] = v
		}
		// 存在过滤键名
		if len(keys) != 0 && !funk.Contains(keys, k) {
			delete(data, k)
		}
	}
	return
}

// Set 分发配置
func (x *Service) Set(ctx context.Context, data map[string]interface{}) (err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		return
	}
	var values map[string]interface{}
	if err = jsoniter.Unmarshal(b, &values); err != nil {
		return
	}
	for k, v := range data {
		values[k] = v
	}
	if b, err = jsoniter.Marshal(values); err != nil {
		return
	}
	if _, err = x.Store.PutBytes("values", b, nats.Context(ctx)); err != nil {
		return
	}
	return
}

// Delete 删除配置
func (x *Service) Delete(ctx context.Context, key string) (err error) {
	var b []byte
	if b, err = x.Store.GetBytes("values", nats.Context(ctx)); err != nil {
		return
	}
	var values map[string]interface{}
	if err = jsoniter.Unmarshal(b, &values); err != nil {
		return
	}
	delete(values, key)
	if b, err = jsoniter.Marshal(values); err != nil {
		return
	}
	if _, err = x.Store.PutBytes("values", b, nats.Context(ctx)); err != nil {
		return
	}
	return
}

// IsSecret 是否为密文
func (x *Service) IsSecret(key string) bool {
	return funk.Contains([]string{
		"tencent_secret_key",
		"tencent_pulsar_token",
		"feishu_app_secret",
		"feishu_encrypt_key",
		"feishu_verification_token",
		"email_password",
		"openapi_secret",
	}, key)
}
