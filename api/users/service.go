package users

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	*common.Inject
}

// FindByIdentity 从用户名或电子邮件获取用户
func (x *Service) FindByIdentity(ctx context.Context, identity string) (data model.User, err error) {
	if err = x.Db.Collection("users").FindOne(ctx, bson.M{
		"status": true,
		"$or": bson.A{
			bson.M{"username": identity},
			bson.M{"email": identity},
		},
	}).Decode(&data); err != nil {
		return
	}
	return
}

// GetActived 获取授权用户数据
func (x *Service) GetActived(ctx context.Context, id string) (data model.User, err error) {
	key := x.Values.Key("users")
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}

	if exists == 0 {
		option := options.Find().SetProjection(bson.M{"password": 0})
		var cursor *mongo.Cursor
		if cursor, err = x.Db.Collection("users").
			Find(ctx, bson.M{"status": true}, option); err != nil {
			return
		}

		values := make(map[string]interface{})
		for cursor.Next(ctx) {
			var user map[string]interface{}
			if err = cursor.Decode(&user); err != nil {
				return
			}

			//var value string
			//if value, err = sonic.MarshalString(user); err != nil {
			//	return
			//}

			values[user["_id"].(string)] = user
		}
		if err = cursor.Err(); err != nil {
			return
		}

		if err = x.Redis.HSet(ctx, key, values).Err(); err != nil {
			return
		}
	}

	if err = x.Redis.HGet(ctx, key, id).Scan(&data); err != nil {
		return
	}

	return
}
