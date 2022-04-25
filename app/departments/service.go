package departments

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindNameById(ctx context.Context, id primitive.ObjectID) (name string, err error) {
	var data common.Department
	if err = x.Db.Collection("departments").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&data); err != nil {
		return
	}
	return data.Name, nil
}
