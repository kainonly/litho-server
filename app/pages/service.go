package pages

import (
	"api/common"
	"api/model"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	*common.Inject
}

func (x *Service) FindOnePage(ctx context.Context, id primitive.ObjectID) (result model.Page, err error) {
	if err = x.Db.Collection("pages").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&result); err != nil {
		return
	}
	return
}

func (x *Service) SchemaKeyExists(ctx context.Context, key string) (code string, err error) {
	var count int64
	if count, err = x.Db.Collection("pages").CountDocuments(ctx, bson.M{
		"schema.key": key,
	}); err != nil {
		return
	}
	if count != 0 {
		return "duplicated", nil
	}
	var colls []string
	if colls, err = x.Db.ListCollectionNames(ctx, bson.M{}); err != nil {
		return
	}
	if funk.Contains(colls, key) {
		return "conflict", nil
	}
	return "", err
}

func (x *Service) Reorganization(ctx context.Context, data ReorganizationDto) (*mongo.BulkWriteResult, error) {
	models := []mongo.WriteModel{
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": data.Id}).
			SetUpdate(bson.M{"$set": bson.M{"parent": data.Parent}}),
	}
	for i, x := range data.Sort {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": x}).
			SetUpdate(bson.M{"$set": bson.M{"sort": i}}),
		)
	}
	return x.Db.Collection("pages").BulkWrite(ctx, models)
}

func (x *Service) SortSchemaFields(ctx context.Context, data SortSchemaFieldsDto) (*mongo.UpdateResult, error) {
	values := make(bson.M, 0)
	for i, x := range data.Fields {
		values[fmt.Sprintf(`schema.fields.%s.sort`, x)] = i
	}
	return x.Db.Collection("pages").UpdateOne(ctx,
		bson.M{"_id": data.Id},
		bson.M{"$set": values},
	)
}

func (x *Service) DeleteSchemaField(ctx context.Context, data DeleteSchemaFieldDto) (*mongo.UpdateResult, error) {
	return x.Db.Collection("pages").UpdateOne(ctx,
		bson.M{"_id": data.Id},
		bson.M{"$unset": bson.M{
			fmt.Sprintf("schema.fields.%s", data.Key): "",
		}},
	)
}

func (x *Service) FindIndexes(ctx context.Context, name string) (result []map[string]interface{}, err error) {
	var cursor *mongo.Cursor
	if cursor, err = x.Db.Collection(name).Indexes().List(ctx); err != nil {
		return
	}
	if err = cursor.All(ctx, &result); err != nil {
		return
	}
	result = result[1:]
	return
}

func (x *Service) CreateIndex(ctx context.Context, data CreateIndexDto, name string) (string, error) {
	return x.Db.Collection(name).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    data.Keys,
		Options: options.Index().SetName(data.Name).SetUnique(*data.Unique),
	})
}

func (x *Service) DeleteIndex(ctx context.Context, data DeleteIndexDto, name string) (bson.Raw, error) {
	return x.Db.Collection(name).Indexes().DropOne(ctx, data.Name)
}

func (x *Service) UpdateValidator(ctx context.Context, data UpdateValidatorDto) (result interface{}, err error) {
	var jsonSchema bson.M
	if err = jsoniter.Unmarshal([]byte(data.Validator), &jsonSchema); err != nil {
		return
	}
	if result, err = x.Db.Collection("pages").UpdateOne(ctx, bson.M{
		"_id": data.Id,
	}, bson.M{
		"$set": bson.M{
			"schema.validator": jsonSchema,
		},
	}); err != nil {
		return
	}
	var page model.Page
	if err = x.Db.Collection("pages").FindOne(ctx, bson.M{
		"_id": data.Id,
	}).Decode(&page); err != nil {
		return
	}
	delete(jsonSchema, "$schema")
	if len(jsonSchema) == 0 {
		return
	}
	if err = x.Db.RunCommand(ctx, bson.D{
		{"collMod", page.Schema.Key},
		{"validator", bson.M{
			"$jsonSchema": jsonSchema,
		}},
	}).Err(); err != nil {
		return
	}
	return
}
