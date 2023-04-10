package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Video struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	ShopId     primitive.ObjectID   `bson:"shop_id" json:"shop_id"`
	Name       string               `bson:"name" json:"name"`
	Url        string               `bson:"url" json:"url"`
	Tags       []primitive.ObjectID `bson:"tags" json:"tags"`
	CreateTime time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime time.Time            `bson:"update_time" json:"update_time"`
}

type VideoTag struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ShopId     primitive.ObjectID `bson:"shop_id" json:"shop_id"`
	Name       string             `bson:"name" json:"name"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}
