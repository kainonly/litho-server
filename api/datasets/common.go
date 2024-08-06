package datasets

import (
	"github.com/google/wire"
	"github.com/weplanx/go/values"
	"github.com/weplanx/server/common"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	DatasetsX *Service
}

type Service struct {
	*common.Inject

	Values *values.Service
}

type M map[string]interface{}

type Dataset struct {
	Name       string   `bson:"name" json:"name"`
	Type       string   `bson:"type" json:"type"`
	Keys       []string `bson:"-" json:"keys"`
	Sensitives []string `bson:"-" json:"sensitives"`
	Status     bool     `bson:"-" json:"status"`
	Event      bool     `bson:"-" json:"event"`
	Options    M        `bson:"options" json:"options"`
}
