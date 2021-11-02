package admin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"laboratory/common"
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

func (x *Service) FindByUsername(ctx context.Context, username string) (data map[string]interface{}, err error) {
	if err = x.Db.Collection("admin").FindOne(ctx, bson.M{
		"username": username,
		"status":   true,
	}).Decode(&data); err != nil {
		return
	}
	return
}
