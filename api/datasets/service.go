package datasets

import (
	"context"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

type Service struct {
	*common.Inject
}

type Dataset struct {
	Name    string   `bson:"name" json:"name"`
	Type    string   `bson:"type" json:"type"`
	Keys    []string `bson:"-" json:"keys"`
	Status  bool     `bson:"-" json:"status"`
	Event   bool     `bson:"-" json:"event"`
	Options M        `bson:"options" json:"options"`
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
		v.Status = control.Status
		v.Event = control.Event
		data = append(data, v)
	}
	return
}
