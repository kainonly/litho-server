package datasets

import (
	"context"
	"github.com/weplanx/go/values"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
)

type Service struct {
	*common.Inject
	Values *values.Service
}

type Dataset struct {
	Name       string   `bson:"name" json:"name"`
	Type       string   `bson:"type" json:"type"`
	Keys       []string `bson:"-" json:"keys"`
	Sensitives []string `bson:"-" json:"sensitives"`
	Status     bool     `bson:"-" json:"status"`
	Event      bool     `bson:"-" json:"event"`
	Options    M        `bson:"options" json:"options"`
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

func (x *Service) Create(ctx context.Context, dto CreateDto) (err error) {
	option := options.CreateCollection()
	if dto.Kind == "timeseries" {
		option = option.
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField(dto.Option.Time).
					SetMetaField(dto.Option.Meta),
			).
			SetExpireAfterSeconds(dto.Option.Expire * 86400)
	}
	if err = x.Db.CreateCollection(ctx, dto.Name, option); err != nil {
		return
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

func (x *Service) Delete(ctx context.Context, name string) (err error) {
	if err = x.Db.Collection(name).Drop(ctx); err != nil {
		return
	}
	controls := x.V.RestControls
	delete(controls, name)
	if err = x.Values.Set(M{
		"RestControls": controls,
	}); err != nil {
		return
	}
	return
}
