package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Picture struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name       string               `bson:"name" json:"name"`
	Url        string               `bson:"url" json:"url"`
	Query      string               `bson:"query" json:"query"`
	Process    PictureProcess       `bson:"process" json:"process"`
	Tags       []primitive.ObjectID `bson:"tags" json:"tags"`
	CreateTime time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime time.Time            `bson:"update_time" json:"update_time"`
}

type PictureProcess struct {
	Mode int64              `bson:"mode" json:"mode"`
	Cut  PictureProcessCut  `bson:"cut" json:"cut"`
	Zoom PictureProcessZoom `bson:"zoom" json:"zoom"`
}

type PictureProcessCut struct {
	X int64 `bson:"x" json:"x"`
	Y int64 `bson:"y" json:"y"`
	W int64 `bson:"w" json:"w"`
	H int64 `bson:"h" json:"h"`
}

type PictureProcessZoom struct {
	W int64 `bson:"w" json:"w"`
	H int64 `bson:"h" json:"h"`
}

type PictureTag struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}
