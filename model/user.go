package model

import (
	"context"
	"github.com/weplanx/go/passlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type User struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email      string               `bson:"email" json:"email"`
	Password   string               `bson:"password" json:"-"`
	Roles      []primitive.ObjectID `bson:"roles" json:"-"`
	Name       string               `bson:"name" json:"name"`
	Avatar     string               `bson:"avatar" json:"avatar"`
	Phone      string               `bson:"phone,omitempty" json:"-"`
	Totp       string               `bson:"totp,omitempty" json:"-"`
	Lark       *UserLark            `bson:"lark,omitempty" json:"lark"`
	Sessions   int64                `bson:"sessions,omitempty" json:"sessions"`
	History    *LogsetLogin         `bson:"history,omitempty" json:"history"`
	Status     bool                 `bson:"status" json:"status"`
	CreateTime time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime time.Time            `bson:"update_time" json:"update_time"`
}

type UserLark struct {
	AccessToken      string `bson:"access_token" json:"access_token"`
	TokenType        string `bson:"token_type" json:"token_type"`
	ExpiresIn        uint64 `bson:"expires_in" json:"expires_in"`
	Name             string `bson:"name" json:"name"`
	EnName           string `bson:"en_name" json:"en_name"`
	AvatarUrl        string `bson:"avatar_url" json:"avatar_url"`
	AvatarThumb      string `bson:"avatar_thumb" json:"avatar_thumb"`
	AvatarMiddle     string `bson:"avatar_middle" json:"avatar_middle"`
	AvatarBig        string `bson:"avatar_big" json:"avatar_big"`
	OpenId           string `bson:"open_id" json:"open_id"`
	UnionId          string `bson:"union_id" json:"union_id"`
	Email            string `bson:"email" json:"email"`
	EnterpriseEmail  string `bson:"enterprise_email" json:"enterprise_email"`
	UserId           string `bson:"user_id" json:"user_id"`
	Mobile           string `bson:"mobile" json:"mobile"`
	TenantKey        string `bson:"tenant_key" json:"tenant_key"`
	RefreshExpiresIn uint64 `bson:"refresh_expires_in" json:"refresh_expires_in"`
	RefreshToken     string `bson:"refresh_token" json:"refresh_token"`
	Sid              string `bson:"sid" json:"sid"`
}

func SetUsers(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "users"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("user", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "users", option); err != nil {
			return
		}
		index := mongo.IndexModel{
			Keys: bson.D{{"email", 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("idx_email"),
		}
		if _, err = db.Collection("users").
			Indexes().
			CreateOne(ctx, index); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "users"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}

func NewUser(email string, password string) *User {
	hash, _ := passlib.Hash(password)
	return &User{
		Email:      email,
		Password:   hash,
		Roles:      []primitive.ObjectID{},
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
