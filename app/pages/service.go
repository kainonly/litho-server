package pages

import (
	"api/common"
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.Inject
}

func (x *Service) FindById(ctx context.Context, id *primitive.ObjectID) (data model.Page, err error) {
	if err = x.Db.Collection("pages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	return
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
