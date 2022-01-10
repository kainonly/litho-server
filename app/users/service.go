package users

import (
	"api/common"
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindByUsername(ctx context.Context, username string) (data model.User, err error) {
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

func (x *Service) FindById(ctx context.Context, id string) (data model.User, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"_id": oid,
		}).
		Decode(&data); err != nil {
		return
	}
	return
}
