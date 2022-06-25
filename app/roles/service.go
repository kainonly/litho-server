package roles

import (
	"api/common"
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOneById(ctx context.Context, id primitive.ObjectID, data interface{}) (err error) {
	if err = x.Db.Collection("roles").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) FindNamesByIds(ctx context.Context, ids []primitive.ObjectID) (names []string, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("roles").
		Find(ctx, bson.M{"_id": bson.M{"$in": ids}}); err != nil {
		return
	}
	var data []model.Role
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	names = make([]string, len(data))
	for k, v := range data {
		names[k] = v.Name
	}
	return
}

func (x *Service) FindByIds(ctx context.Context, ids []primitive.ObjectID, data interface{}) (err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("roles").
		Find(ctx, bson.M{"_id": bson.M{"$in": ids}}); err != nil {
		return
	}
	if err = cursor.All(ctx, data); err != nil {
		return
	}
	return
}
