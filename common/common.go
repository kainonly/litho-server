package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/weplanx/support/helper"
	"github.com/weplanx/support/passport"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

type App struct {
	fx.In

	Set    *Set
	Mongo  *mongo.Client
	Db     *mongo.Database
	Redis  *redis.Client
	Cookie *helper.CookieHelper
	Cipher *helper.CipherHelper
}

type Set struct {
	Name     string                        `yaml:"name"`
	Key      string                        `yaml:"key"`
	Database Database                      `yaml:"database"`
	Redis    Redis                         `yaml:"redis"`
	Cookie   helper.CookieOption           `yaml:"cookie"`
	Cors     []string                      `yaml:"cors"`
	Auth     map[string]*passport.Passport `yaml:"auth"`
}

type Database struct {
	Uri  string `yaml:"uri"`
	Name string `yaml:"name"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (x *Set) RedisKey(name string) string {
	return x.Name + ":" + name
}
