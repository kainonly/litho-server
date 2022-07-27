package pages

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	*common.Inject
}

type NavDto struct {
	ID     primitive.ObjectID `json:"_id"`
	Parent interface{}        `json:"parent"`
	Name   string             `json:"name"`
	Icon   string             `json:"icon"`
	Kind   string             `json:"kind"`
	Sort   int64              `json:"sort"`
}

func (x *Service) Navs(ctx context.Context, roles []model.Role) (navs []NavDto, err error) {
	pageIds := make([]primitive.ObjectID, 0)
	pageSet := make(map[string]bool)
	for _, role := range roles {
		for k := range role.Pages {
			if pageSet[k] {
				continue
			}
			id, _ := primitive.ObjectIDFromHex(k)
			pageIds = append(pageIds, id)
			pageSet[k] = true
		}
	}
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("pages").
		Find(ctx, bson.M{
			"_id":    bson.M{"$in": pageIds},
			"status": true,
		}); err != nil {
		return
	}
	if err = cursor.All(ctx, &navs); err != nil {
		return
	}
	return
}
