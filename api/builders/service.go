package builders

import (
	"context"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*common.Inject
}

func (x *Service) SortFields(ctx context.Context, id primitive.ObjectID, keys []string) (err error) {
	var builder model.Builder
	if err = x.Db.Collection("builders").
		FindOne(ctx, M{"_id": id}).
		Decode(&builder); err != nil {
		return
	}
	dict := make(map[string]model.BuilderSchemaField)
	for _, v := range builder.Schema.Fields {
		dict[v.Key] = v
	}
	data := make([]model.BuilderSchemaField, len(dict))
	for i, key := range keys {
		data[i] = dict[key]
	}
	update := M{"$set": M{"schema.fields": data}}
	if _, err = x.Db.Collection("builders").
		UpdateByID(ctx, id, update); err != nil {
		return
	}
	return
}
