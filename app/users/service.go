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
