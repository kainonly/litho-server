package pages

import (
	"context"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	*common.Inject
}

type Nav struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Parent interface{}        `json:"parent"`
	Name   string             `json:"name"`
	Icon   string             `json:"icon"`
	Kind   string             `json:"kind"`
	Sort   int64              `json:"sort"`
}

// FindNavs 筛选导航数据
func (x *Service) FindNavs(ctx context.Context) (data []Nav, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("pages").
		Find(ctx, bson.M{"status": true}); err != nil {
		return
	}
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}
