package users

import (
	"api/common"
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	*common.Inject
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

func (x *Service) FindLabels(ctx context.Context) (values []interface{}, err error) {
	if values, err = x.Db.Collection("users").
		Distinct(ctx, "labels", bson.M{"status": true}); err != nil {
		return
	}
	return
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

func (x *Service) FindInfo(ctx context.Context, id string) (data model.User, err error) {
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	if err = x.Db.Collection("users").FindOne(ctx,
		bson.M{"_id": oid},
		options.FindOne().SetProjection(bson.M{
			"password": 0,
			"roles":    0,
			"pages":    0,
			"readonly": 0,
		}),
	).Decode(&data); err != nil {
		return
	}
	return
}
