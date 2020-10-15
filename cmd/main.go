package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"os"
	"time"
	"van-api/app/controller"
	"van-api/app/middleware/cors"
	"van-api/app/types"
)

func main() {
	if _, err := os.Stat("./config/config.yml"); os.IsNotExist(err) {
		log.Fatalln("the configuration file does not exist")
	}
	buf, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("failed to read service configuration file", err)
	}
	cfg := types.Config{}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		log.Fatalln("service configuration file parsing failed", err)
	}
	app := iris.Default()
	app.Use(cors.Cors(types.CorsOption{}))
	context.DefaultJSONOptions = context.JSON{
		StreamingJSON: true,
	}
	db, err := gorm.Open(mysql.Open(cfg.Mysql.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Mysql.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
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
	app.ConfigureContainer(func(container *router.APIContainer) {
		container.RegisterDependency(db)
		container.RegisterDependency(rdb)
		container.Get("/", controller.Default)
		container.Options("*", controller.Default)
		main := container.Party("/main")
		{
			main.Post("/verify", controller.MainVerify)
		}
	})
	app.Listen(cfg.Listen, iris.WithOptimizations)
}
