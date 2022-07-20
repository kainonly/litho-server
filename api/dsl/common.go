package dsl

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type M = map[string]interface{}

type FindOption struct {
	// 排序规则
	Sort M
	// 投影规则
	Keys M
	// 返回数量
	Limit int64
	// 跳过数量
	Skip int64
	// 页码
	Page int64
}

func (x *FindOption) GetSort() (data bson.D) {
	for key, value := range x.Sort {
		data = append(data, bson.E{Key: key, Value: value})
	}
	return
}
