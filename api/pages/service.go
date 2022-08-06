package pages

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetNavs 筛选导航数据
func (x *Service) GetNavs(ctx context.Context) (data []Nav, err error) {
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

// FindOneById 通过 ID 查找
func (x *Service) FindOneById(ctx context.Context, id primitive.ObjectID) (data model.Page, err error) {
	if err = x.Db.Collection("pages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&data); err != nil {
		return
	}
	return
}

// GetIndexes 获取索引
func (x *Service) GetIndexes(ctx context.Context, name string) (indexes []map[string]interface{}, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(name).
		Indexes().
		List(ctx); err != nil {
		return
	}
	indexes = make([]map[string]interface{}, 0)
	if err = cursor.All(ctx, &indexes); err != nil {
		return
	}
	return
}

// SetIndex 设置索引
func (x *Service) SetIndex(ctx context.Context, coll string, name string, keys bson.D, unique bool) (string, error) {
	return x.Db.Collection(coll).
		Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys: keys,
			Options: options.Index().
				SetName(name).
				SetUnique(unique),
		})
}

// DeleteIndex 删除索引
func (x *Service) DeleteIndex(ctx context.Context, coll string, name string) (bson.Raw, error) {
	return x.Db.Collection(coll).
		Indexes().
		DropOne(ctx, name)
}
