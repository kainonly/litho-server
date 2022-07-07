package departments

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"server/common"
	"server/model"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOneById(ctx context.Context, id primitive.ObjectID, data interface{}) (err error) {
	if err = x.Db.Collection("departments").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(data); err != nil {
		return
	}
	return
}

func (x *Service) FindNameById(ctx context.Context, id primitive.ObjectID) (name string, err error) {
	var data model.Department
	if err = x.FindOneById(ctx, id, &data); err != nil {
		return
	}
	return data.Name, nil
}
