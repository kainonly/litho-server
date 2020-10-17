package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"os"
	"time"
	"van-api/app"
	"van-api/helper"
	"van-api/route"
	"van-api/types"
)

func main() {
	if _, err := os.Stat("./config/config.yml"); os.IsNotExist(err) {
		log.Fatalln("the configuration file does not exist")
	}
	buf, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("failed to read service configuration file", err)
	}
	var cfg types.Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		log.Fatalln("service configuration file parsing failed", err)
	}
	serve := iris.Default()
	serve.Use(cors.New(cors.Options{
		AllowedOrigins:   cfg.Cors.Origin,
		AllowedMethods:   cfg.Cors.Method,
		AllowedHeaders:   cfg.Cors.AllowHeader,
		ExposedHeaders:   cfg.Cors.ExposedHeader,
		MaxAge:           cfg.Cors.MaxAge,
		AllowCredentials: cfg.Cors.Credentials,
	}))

	db, err := gorm.Open(mysql.Open(cfg.Mysql.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Mysql.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	helper.DB = db
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	if cfg.Mysql.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
	}
	if cfg.Mysql.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
	}
	if cfg.Mysql.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.Mysql.ConnMaxLifetime))
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	helper.RDB = rdb
	serve.ConfigureContainer(func(container *router.APIContainer) {
		container.RegisterDependency(db)
		container.RegisterDependency(rdb)
		container.Get("/", route.Default)
		container.Options("*", route.Default)
		mvc.Configure(container.Party("/").Self, app.Bootstrap)
	})
	serve.Listen(cfg.Listen, iris.WithOptimizations)
}
