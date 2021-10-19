package page

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"laboratory/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}

func (x *Service) Get(ctx context.Context) (data []map[string]interface{}, err error) {
	opt := options.Find()
	opt.Projection = bson.M{
		"option": 0,
	}
	cursor, err := x.Db.Collection("page").Find(ctx, bson.M{}, opt)
	if err != nil {
		return
	}
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}
