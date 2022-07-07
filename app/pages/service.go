package pages

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/common"
	"server/model"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOneById(ctx context.Context, id primitive.ObjectID, data interface{}) (err error) {
	if err = x.Db.Collection("pages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(data); err != nil {
		return
	}
	return
}

type NavDto struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Parent interface{}        `bson:"parent" json:"parent"`
	Name   string             `bson:"name" json:"name"`
	Icon   string             `bson:"icon" json:"icon"`
	Kind   string             `bson:"kind" json:"kind"`
	Sort   int64              `bson:"sort" json:"sort"`
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
