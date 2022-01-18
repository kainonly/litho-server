package users

import (
	"api/common"
	"api/model"
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindLabels(ctx context.Context) (values []model.Value, err error) {
	var result []interface{}
	if result, err = x.Db.Collection("users").
		Distinct(ctx, "labels", bson.M{"status": true}); err != nil {
		return
	}
	values = make([]model.Value, 0)
	if len(result) == 0 {
		return
	}
	for _, data := range result {
		var value model.Value
		for _, v := range data.(primitive.D) {
			if v.Key == "label" {
				value.Label = v.Value.(string)
			}
			if v.Key == "value" {
				value.Value = v.Value
			}
		}
		values = append(values, value)
	}
	return
}

func (x *Service) HasUsername(ctx context.Context, username string) (code string, err error) {
	var count int64
	if count, err = x.Db.Collection("users").CountDocuments(ctx, bson.M{
		"username": username,
	}); err != nil {
		return
	}
	if count != 0 {
		return "duplicated", nil
	}
	return "", err
}

func (x *Service) FindOneByUsername(ctx context.Context, username string) (data model.User, err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"username": username,
			"status":   true,
		}).
		Decode(&data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneFromCacheById(ctx context.Context, id string) (data model.User, err error) {
	key := x.Values.KeyName("users", id)
	var value []byte
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if exists != 0 {
		if err = x.Redis.Get(ctx, key).Scan(&value); err != nil {
			return
		}
		if err = jsoniter.Unmarshal(value, &data); err != nil {
			return
		}
		return
	}
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	if err = x.Db.Collection("users").
		FindOne(ctx,
			bson.M{"_id": oid},
			options.FindOne().SetProjection(bson.M{
				"password": 0,
				"status":   0,
			}),
		).
		Decode(&data); err != nil {
		return
	}
	if value, err = jsoniter.Marshal(data); err != nil {
		return
	}
	if err = x.Redis.Set(ctx, key, value, 0).Err(); err != nil {
		return
	}
	return
}
