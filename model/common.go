package model

import "go.mongodb.org/mongo-driver/bson/primitive"

func Bool(v bool) *bool {
	return &v
}

func ObjectID(v interface{}) *primitive.ObjectID {
	if id, ok := v.(primitive.ObjectID); ok {
		return &id
	}
	return nil
}

type Value struct {
	// 名称
	Label string `bson:"label" json:"label"`

	// 数值
	Value interface{} `bson:"value" json:"value"`
}
