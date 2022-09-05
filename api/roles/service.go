package roles

import (
	"context"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Db *mongo.Database
}

func (x *Service) FindByIds(ctx context.Context, ids []primitive.ObjectID) (err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("roles").
		Find(ctx, bson.M{"_id": bson.M{"$in": ids}}); err != nil {
		return
	}

	//pages := make(map[string]*int64)
	for cursor.Next(ctx) {
		var value model.Role
		if err = cursor.Decode(&value); err != nil {
			return
		}

		// TODO: 权限待改造
	}

	if err = cursor.Err(); err != nil {
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

	names = make([]string, 0)
	for cursor.Next(ctx) {
		var value model.Role
		if err = cursor.Decode(&value); err != nil {
			return
		}

		names = append(names, value.Name)
	}

	if err = cursor.Err(); err != nil {
		return
	}

	return
}
