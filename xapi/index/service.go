package index

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	*common.Inject
}

func (x *Service) Accelerate(ctx context.Context) (result []model.AccTask, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("acc_tasks").Find(ctx, bson.M{}); err != nil {
		return
	}
	if err = cursor.All(ctx, &result); err != nil {
		return
	}
	return
}
