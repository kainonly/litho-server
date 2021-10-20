package devops

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"laboratory/common"
)

type Service struct {
	*InjectService
}

type InjectService struct {
	common.App
}

func (x *Service) InitData(ctx context.Context) (err error) {
	if _, err = x.Db.Collection("role").InsertOne(ctx, bson.M{
		"key":         "*",
		"name":        "超级管理员",
		"description": "",
		"pages":       bson.D{},
	}); err != nil {
		return
	}
	if _, err = x.Db.Collection("admin").InsertOne(ctx, bson.M{
		"username": "admin",
		"password": "",
		"roles":    bson.A{"*"},
		"name":     "",
		"email":    "",
		"phone":    "",
		"avatar":   "",
	}); err != nil {
		return
	}
	return
}
