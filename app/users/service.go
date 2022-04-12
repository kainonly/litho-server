package users

import (
	"api/common"
	"api/common/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOneByUsernameOrEmail(ctx context.Context, value string) (data model.User, err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"status": true,
			"$or": bson.A{
				bson.M{"username": value},
				bson.M{"email": value},
			},
		}).
		Decode(&data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneById(ctx context.Context, id string, data interface{}, opts ...*options.FindOneOptions) (err error) {
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	if err = x.Db.Collection("users").FindOne(ctx,
		bson.M{"_id": oid},
		opts...,
	).Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) UpdateOneById(ctx context.Context, id string, update interface{}) (err error) {
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	if _, err = x.Db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": oid},
		update,
	); err != nil {
		return
	}
	return
}
