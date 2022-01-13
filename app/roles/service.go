package roles

import (
	"api/common"
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindLabels(ctx context.Context) (values []model.Value, err error) {
	var result []interface{}
	if result, err = x.Db.Collection("roles").
		Distinct(ctx, "labels", bson.M{"status": true}); err != nil {
		return
	}
	values = make([]model.Value, 0)
	if len(result) == 0 {
		return
	}
	for _, data := range result {
		var value model.Value
		for _, v := range data.(primitive.D) {
			if v.Key == "label" {
				value.Label = v.Value.(string)
			}
			if v.Key == "value" {
				value.Value = v.Value
			}
		}
		values = append(values, value)
	}
	return
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
