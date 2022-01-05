package users

import (
	"api/common"
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindByUsername(ctx context.Context, username string) (data model.User, err error) {
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"username": username,
			"status":   true,
		}).
		Decode(&data); err != nil {
		return
	}
	return
}
