package datasets

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/values"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateDto struct {
	Name   string           `json:"name" vd:"required"`
	Kind   string           `json:"kind" vd:"oneof='collection' 'timeseries'"`
	Option *CreateOptionDto `json:"option" vd:"required_if=kind 'timeseries'"`
}

type CreateOptionDto struct {
	Time   string `json:"time" vd:"required"`
	Meta   string `json:"meta" vd:"required"`
	Expire *int64 `json:"expire" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.DatasetsX.Create(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func (x *Service) Create(ctx context.Context, dto CreateDto) (err error) {
	var names []string
	if names, err = x.Db.ListCollectionNames(ctx, bson.M{"name": dto.Name}); err != nil {
		return
	}
	if len(names) == 0 {
		option := options.CreateCollection()
		if dto.Kind == "timeseries" {
			option = option.
				SetTimeSeriesOptions(
					options.TimeSeries().
						SetTimeField(dto.Option.Time).
						SetMetaField(dto.Option.Meta),
				).
				SetExpireAfterSeconds(*dto.Option.Expire * 86400)
		}

		if err = x.Db.CreateCollection(ctx, dto.Name, option); err != nil {
			return
		}
	}
	controls := x.V.RestControls
	controls[dto.Name] = &values.RestControl{
		Keys:   nil,
		Status: true,
		Event:  false,
	}
	if err = x.Values.Set(M{
		"RestControls": controls,
	}); err != nil {
		return
	}
	return
}
