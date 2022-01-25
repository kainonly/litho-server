package videos

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindLabels(ctx context.Context) (values []interface{}, err error) {
	if values, err = x.Db.Collection("media").
		Distinct(ctx, "labels", bson.M{"status": true}); err != nil {
		return
	}
	return
}

func (x *Service) BulkDelete(ctx context.Context, oids []primitive.ObjectID) (interface{}, error) {
	return x.Db.Collection("media").DeleteMany(ctx, bson.M{
		"_id": bson.M{"$in": oids},
	})
}
