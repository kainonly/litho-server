package bootstrap

import (
	"context"
	"errors"
	"github.com/weplanx/server/model"
	"github.com/weplanx/server/utils/passlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Install struct {
	Db       *mongo.Database
	Username string
	Password string
}

func (x *Install) Basic(ctx context.Context) (err error) {
	var exists []string
	if exists, err = x.Db.ListCollectionNames(ctx, bson.M{
		"name": bson.M{"$in": bson.A{"roles", "users"}},
	}); err != nil {
		return
	}
	if len(exists) != 0 {
		return errors.New("操作不被允许, [roles] 与 [users] 集合是存在的")
	}
	// 初始化权限组
	var roles *mongo.InsertOneResult
	if roles, err = x.Db.Collection("roles").
		InsertOne(ctx,
			model.NewRole("超级管理员").
				SetDescription("系统默认设置"),
		); err != nil {
		return
	}
	if _, err = x.Db.Collection("roles").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"name": 1},
				Options: options.Index().SetName("uk_name").SetUnique(true),
			},
		},
	); err != nil {
		return
	}

	// 初始化管理用户
	passwordHash, _ := passlib.Hash(x.Password)
	user := model.NewUser(x.Username, passwordHash).
		SetRoles([]primitive.ObjectID{roles.InsertedID.(primitive.ObjectID)})
	if _, err = x.Db.Collection("users").
		InsertOne(ctx, user); err != nil {
		return
	}
	if _, err = x.Db.Collection("users").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"username": 1},
				Options: options.Index().SetName("idx_username").SetUnique(true),
			},
			{
				Keys:    bson.M{"email": 1},
				Options: options.Index().SetName("idx_email"),
			},
		},
	); err != nil {
		return
	}

	return
}
