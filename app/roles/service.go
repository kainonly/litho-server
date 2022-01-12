package roles

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindLabels(ctx context.Context) ([]interface{}, error) {
	return x.Db.Collection("roles").
		Distinct(ctx, "labels", bson.M{"status": true})
}

func (x *Service) HasKey(ctx context.Context, key string) (code string, err error) {
	var count int64
	if count, err = x.Db.Collection("roles").CountDocuments(ctx, bson.M{
		"key": key,
	}); err != nil {
		return
	}
	if count != 0 {
		return "duplicated", nil
	}
	return "", err
}
