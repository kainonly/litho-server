package media

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
