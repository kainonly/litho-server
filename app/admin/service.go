package admin

import (
	"api/common"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.Inject
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
