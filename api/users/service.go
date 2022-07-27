package users

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Actived 获取授权用户数据
func (x *Service) Actived(ctx context.Context, id primitive.ObjectID) (data map[string]interface{}, err error) {
	// TODO: 更换为 Redis
	option := options.FindOne().
		SetProjection(bson.M{"password": 0})
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id, "status": true}, option).
		Decode(&data); err != nil {
		return
	}
	return
}
