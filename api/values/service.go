package values

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/api/common"
	"github.com/weplanx/api/model"
	"time"
)

type Service struct {
	*common.Inject
}

// Key 命名
func (x *Service) Key() string {
	return fmt.Sprintf("%s:values", x.Values.App.Namespace)
}

// Load 载入配置
func (x *Service) Load(ctx context.Context) (err error) {
	var count int64
	if count, err = x.Redis.Exists(ctx, x.Key()).Result(); err != nil {
		return
	}
	var b []byte
	// 不存在配置则初始化
	if count == 0 {
		x.Values.DynamicValues = common.DynamicValues{
			"session_ttl":               float64(3600),
			"login_ttl":                 float64(900),
			"login_failures":            float64(5),
			"ip_login_failures":         float64(10),
			"ip_whitelist":              []string{},
			"ip_blacklist":              []string{},
			"pwd_strategy":              float64(1),
			"pwd_ttl":                   float64(365),
			"cloud":                     "",
			"tencent_secret_id":         "",
			"tencent_secret_key":        "",
			"tencent_cos_bucket":        "",
			"tencent_cos_region":        "",
			"tencent_cos_expired":       float64(300),
			"tencent_cos_limit":         float64(5120),
			"office":                    "",
			"feishu_app_id":             "",
			"feishu_app_secret":         "",
			"feishu_encrypt_key":        "",
			"feishu_verification_token": "",
			"redirect_url":              "",
			"email_host":                "",
			"email_port":                "465",
			"email_username":            "",
			"email_password":            "",
			"openapi_url":               "",
			"openapi_key":               "",
			"openapi_secret":            "",
		}

		if b, err = sonic.Marshal(x.Values.DynamicValues); err != nil {
			return
		}

		if err = x.Redis.Set(ctx, x.Key(), b, 0).Err(); err != nil {
			return
		}

		return
	}

	if b, err = x.Redis.Get(ctx, x.Key()).Bytes(); err != nil {
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

	if _, err = x.Nats.Subscribe(x.Key(), func(msg *nats.Msg) {
		if string(msg.Data) != "sync" {
			return
		}
		if err = x.Load(context.TODO()); err != nil {
			fmt.Println(err)
		}
	}); err != nil {
		return
	}

	return
}

// Get 获取动态配置
func (x *Service) Get(keys ...string) (data map[string]interface{}) {
	sets := make(map[string]bool)
	for _, key := range keys {
		sets[key] = true
	}
	isAll := len(sets) == 0
	data = make(map[string]interface{})
	for k, v := range x.Values.DynamicValues {
		if !isAll && !sets[k] {
			continue
		}
		if Secret[k] {
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
	// 合并覆盖
	for k, v := range data {
		x.Values.DynamicValues[k] = v
	}
	return x.Update(ctx)
}

// Remove 移除动态配置
func (x *Service) Remove(ctx context.Context, key string) (err error) {
	delete(x.Values.DynamicValues, key)
	return x.Update(ctx)
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

	// 写入日志
	if _, err = x.Db.Collection("values_logs").
		InsertOne(ctx, model.ValuesLog{
			Time:     time.Now(),
			Uid:      ctx.Value("uid"),
			Snapshot: x.Values.DynamicValues,
		}); err != nil {
		return
	}
	return
}
