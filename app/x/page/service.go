package page

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type InjectService struct {
	*common.App
}

type Service struct {
	*InjectService
}

func NewService(i *InjectService) *Service {
	return &Service{
		InjectService: i,
	}
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
