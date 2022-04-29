package vars

import (
	"api/common"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

type Service struct {
	*common.Inject
}

// Get 获取变量
func (x *Service) Get(ctx context.Context, key string) (value string, err error) {
	if err = x.Refresh(ctx); err != nil {
		return
	}
	return x.Redis.HGet(ctx, x.Values.KeyName("vars"), key).Result()
}

// Gets 获取指定变量
func (x *Service) Gets(ctx context.Context, keys []string) (values map[string]interface{}, err error) {
	if err = x.Refresh(ctx); err != nil {
		return
	}
	var result []interface{}
	if result, err = x.Redis.HMGet(ctx, x.Values.KeyName("vars"), keys...).Result(); err != nil {
		return
	}
	values = make(map[string]interface{})
	for k, v := range keys {
		values[v] = result[k]
	}
	return
}

// Refresh 刷新变量
func (x *Service) Refresh(ctx context.Context) (err error) {
	key := x.Values.KeyName("vars")
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var cursor *mongo.Cursor
		if cursor, err = x.Db.Collection("vars").Find(ctx, bson.M{}); err != nil {
			return
		}
		var data []common.Var
		if err = cursor.All(ctx, &data); err != nil {
			return
		}
		pipe := x.Redis.Pipeline()
		for _, v := range data {
			switch x := v.Value.(type) {
			case primitive.A:
				b, _ := jsoniter.Marshal(x)
				pipe.HSet(ctx, key, v.Key, b)
				break
			case primitive.M:
				b, _ := jsoniter.Marshal(x)
				pipe.HSet(ctx, key, v.Key, b)
				break
			default:
				pipe.HSet(ctx, key, v.Key, x)
			}
		}
		if _, err = pipe.Exec(ctx); err != nil {
			return
		}
	}
	return
}

// Set 设置变量
func (x *Service) Set(ctx context.Context, key string, value interface{}) (err error) {
	var exists int64
	if exists, err = x.Db.Collection("vars").CountDocuments(ctx, bson.M{"key": key}); err != nil {
		return
	}
	doc := common.NewVar(key, value)
	if exists == 0 {
		if _, err = x.Db.Collection("vars").InsertOne(ctx, doc); err != nil {
			return
		}
	} else {
		if _, err = x.Db.Collection("vars").ReplaceOne(ctx, bson.M{"key": key}, doc); err != nil {
			return
		}
	}
	if err = x.Redis.Del(ctx, x.Values.KeyName("vars")).Err(); err != nil {
		return
	}
	return
}

type UploadDto struct {
	Type  string `json:"type"`
	Url   string `json:"url"`
	Limit int    `json:"limit"`
}

func (x *Service) GetUpload(ctx context.Context) (data *UploadDto, err error) {
	var platform string
	if platform, err = x.Get(ctx, "cloud_platform"); err != nil {
		return
	}
	switch platform {
	case "tencent":
		var option map[string]interface{}
		if option, err = x.Gets(ctx, []string{
			"tencent_cos_bucket",
			"tencent_cos_region",
			"tencent_cos_limit",
		}); err != nil {
			return
		}
		limit, _ := strconv.Atoi(option["tencent_cos_limit"].(string))
		data = &UploadDto{
			Type: "cos",
			Url: fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
				option["tencent_cos_bucket"], option["tencent_cos_region"],
			),
			Limit: limit,
		}
		break
	}
	return
}

// GetExpiration 获取会话有效时间
func (x *Service) GetExpiration(ctx context.Context) (t time.Duration) {
	value, _ := x.Get(ctx, "user_session_expire")
	if value != "" {
		t, _ = time.ParseDuration(value)
	} else {
		t = time.Hour
	}
	return
}
