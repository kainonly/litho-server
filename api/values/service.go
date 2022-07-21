package values

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	Values *common.Values
	Db     *mongo.Database
	Redis  *redis.Client
	Nats   *nats.Conn
}

// Key 命名
func (x *Service) Key() string {
	return fmt.Sprintf("%s:values", x.Values.App.Namespace)
}

// Load 载入配置
func (x *Service) Load(ctx context.Context) (err error) {
	var b []byte
	if b, err = x.Redis.Get(ctx, x.Key()).Bytes(); err != nil {
		return
	}

	if err = json.Unmarshal(b, &x.Values.DynamicValues); err != nil {
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
	if b, err = json.Marshal(x.Values.DynamicValues); err != nil {
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
