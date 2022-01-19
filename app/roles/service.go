package roles

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	*common.Inject
}

func (x *Service) HasName(ctx context.Context, name string) (code string, err error) {
	var count int64
	if count, err = x.Db.Collection("roles").CountDocuments(ctx, bson.M{
		"name": name,
	}); err != nil {
		return
	}
	if count != 0 {
		return "duplicated", nil
	}
	return "", err
}

func (x *Service) FindLabels(ctx context.Context) (values []interface{}, err error) {
	if values, err = x.Db.Collection("roles").
		Distinct(ctx, "labels", bson.M{"status": true}); err != nil {
		return
	}
	return
}
