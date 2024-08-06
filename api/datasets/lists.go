package datasets

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

type ListsDto struct {
	Name string `query:"name"`
}

func (x *Controller) Lists(ctx context.Context, c *app.RequestContext) {
	var dto ListsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.DatasetsX.Lists(ctx, dto.Name)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) Lists(ctx context.Context, name string) (data []Dataset, err error) {
	var names []string
	for key, _ := range x.V.RestControls {
		var match bool
		if match, err = regexp.Match("^"+name, []byte(key)); err != nil {
			return
		}
		if match {
			names = append(names, key)
		}
	}
	var cursor *mongo.Cursor
	if cursor, err = x.Db.ListCollections(ctx,
		bson.M{"name": bson.M{"$in": names}},
	); err != nil {
		return
	}

	for cursor.Next(ctx) {
		var v Dataset
		if err = cursor.Decode(&v); err != nil {
			return
		}
		control := x.V.RestControls[v.Name]
		v.Keys = control.Keys
		v.Sensitives = control.Sensitives
		v.Status = control.Status
		v.Event = control.Event
		data = append(data, v)
	}
	return
}
