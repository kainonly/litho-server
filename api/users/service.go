package users

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/server/api/departments"
	"github.com/weplanx/server/api/roles"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	Values *common.Values
	Db     *mongo.Database
	Redis  *redis.Client

	RolesService       *roles.Service
	DepartmentsService *departments.Service
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
	key := x.Values.Name("users")
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

		values := make(map[string]string)
		for cursor.Next(ctx) {
			var user model.User
			if err = cursor.Decode(&user); err != nil {
				return
			}

			var value string
			if value, err = sonic.MarshalString(user); err != nil {
				return
			}

			values[user.ID.Hex()] = value
		}
		if err = cursor.Err(); err != nil {
			return
		}

		if err = x.Redis.HSet(ctx, key, values).Err(); err != nil {
			return
		}
	}

	var result string
	if result, err = x.Redis.HGet(ctx, key, id).Result(); err != nil {
		return
	}
	if err = sonic.UnmarshalString(result, &data); err != nil {
		return
	}

	return
}

// UpdateOneById 通过 ID 更新
func (x *Service) UpdateOneById(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	return x.Db.Collection("users").UpdateOne(ctx, bson.M{"_id": id}, update)
}
