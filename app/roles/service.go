package roles

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindNamesById(ctx context.Context, ids []primitive.ObjectID) (names []string, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("roles").Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	}); err != nil {
		return
	}
	var data []common.Role
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	names = make([]string, len(data))
	for k, v := range data {
		names[k] = v.Name
	}
	return
}
