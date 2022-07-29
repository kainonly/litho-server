package departments

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOneById(ctx context.Context, id primitive.ObjectID) (data model.Department, err error) {
	if err = x.Db.Collection("departments").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	return
}
