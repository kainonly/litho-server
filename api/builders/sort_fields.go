package builders

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SortFieldsDto struct {
	Id   primitive.ObjectID `json:"id" vd:"required"`
	Keys []string           `json:"keys" vd:"gt=0"`
}

func (x *Controller) SortFields(ctx context.Context, c *app.RequestContext) {
	var dto SortFieldsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.BuildersX.SortFields(ctx, dto.Id, dto.Keys); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
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
