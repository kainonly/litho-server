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

func (x *Service) FindOneById(ctx context.Context, id string, opts ...*options.FindOneOptions) (data model.User, err error) {
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return
	}
	if err = x.Db.Collection("users").FindOne(ctx,
		bson.M{"_id": oid},
		opts...,
	).Decode(&data); err != nil {
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
