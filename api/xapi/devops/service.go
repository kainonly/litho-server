package devops

import (
	"context"
	"github.com/alexedwards/argon2id"
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
		"status":      true,
		"description": "",
		"pages":       bson.A{},
	}); err != nil {
		return
	}
	var password string
	if password, err = argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams); err != nil {
		return
	}
	if _, err = x.Db.Collection("admin").InsertOne(ctx, bson.M{
		"username": "admin",
		"password": password,
		"status":   true,
		"roles":    bson.A{"*"},
		"name":     "超级管理员",
		"email":    "",
		"phone":    "",
		"avatar":   "",
	}); err != nil {
		return
	}
	return
}
