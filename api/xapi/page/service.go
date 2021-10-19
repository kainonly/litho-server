package page

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"laboratory/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}

func (x *Service) Get(ctx context.Context) (data []map[string]interface{}, err error) {
	cursor, err := x.Db.Collection("page").Find(ctx, bson.M{})
	if err != nil {
		return
	}
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}
