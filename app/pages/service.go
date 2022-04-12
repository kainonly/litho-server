package pages

import (
	"api/common"
	"api/common/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	*common.Inject
}

type NavDto struct {
	ID     primitive.ObjectID `bson:"_id" json:"_id"`
	Parent interface{}        `bson:"parent" json:"parent"`
	Name   string             `bson:"name" json:"name"`
	Icon   string             `bson:"icon" json:"icon"`
	Kind   string             `bson:"kind" json:"kind"`
	Sort   int64              `bson:"sort" json:"sort"`
}

func (x *Service) Navs(ctx context.Context) (data []NavDto, err error) {
	data = make([]NavDto, 0)
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection("pages").Find(ctx,
		bson.M{"status": true},
		options.Find().SetProjection(bson.M{
			"schema":      0,
			"status":      0,
			"create_time": 0,
			"update_time": 0,
		}),
	); err != nil {
		return
	}
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}

func (x *Service) FindOneById(ctx context.Context, id string) (result model.Page, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	if err = x.Db.Collection("pages").
		FindOne(ctx, bson.M{"_id": oid}).
		Decode(&result); err != nil {
		return
	}
	return
}

func (x *Service) Indexes(ctx context.Context, name string) (data []map[string]interface{}, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(name).Indexes().
		List(ctx); err != nil {
		return
	}
	data = make([]map[string]interface{}, 0)
	if err = cursor.All(ctx, &data); err != nil {
		return
	}
	return
}

func (x *Service) CreateIndex(ctx context.Context, coll string, name string, keys bson.D, unique bool) (string, error) {
	return x.Db.Collection(coll).
		Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys: keys,
			Options: options.Index().
				SetName(name).
				SetUnique(unique),
		})
}

func (x *Service) DeleteIndex(ctx context.Context, coll string, name string) (bson.Raw, error) {
	return x.Db.Collection(coll).Indexes().DropOne(ctx, name)
}

func (x *Service) UpdateValidator(ctx context.Context, coll string, validator string) error {
	return x.Db.RunCommand(ctx, bson.D{
		{"collMod", coll},
		{"validator", bson.M{"$jsonSchema": validator}},
	}).Err()
}
