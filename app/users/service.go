package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/common"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOneById(ctx context.Context, id primitive.ObjectID, data interface{}, opts ...*options.FindOneOptions) (err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id}, opts...).
		Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneByUsernameOrEmail(ctx context.Context, search string, data interface{}) (err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"status": true,
			"$or": bson.A{
				bson.M{"username": search},
				bson.M{"email": search},
			},
		}).
		Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneByEmail(ctx context.Context, email string, data interface{}) (err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"email":  email,
			"status": true,
		}).
		Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneByFeishu(ctx context.Context, openid string, data interface{}) (err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"status":        true,
			"feishu.openid": openid,
		}).
		Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) UpdateOneById(ctx context.Context, id primitive.ObjectID, update interface{}) (err error) {
	if _, err = x.Db.Collection("users").
		UpdateOne(ctx, bson.M{"_id": id}, update); err != nil {
		return
	}
	return
}

func (x *Service) UpdateOneByEmail(ctx context.Context, email string, update interface{}) (err error) {
	if _, err = x.Db.Collection("users").
		UpdateOne(ctx, bson.M{"email": email}, update); err != nil {
		return
	}
	return
}

func (x *Service) Count(ctx context.Context, filter bson.M) (count int64, err error) {
	return x.Db.Collection("users").CountDocuments(ctx, filter)
}
