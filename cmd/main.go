package main

import (
	"github.com/iris-contrib/middleware/cors"
	curd "github.com/kainonly/iris-curd"
	"github.com/kainonly/iris-helper/token"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"van-api/app"
	"van-api/app/cache"
	"van-api/bootstrap"
	"van-api/route"
)

func main() {
	cfg, err := bootstrap.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	serve := iris.New()
	serve.Use(recover.New())
	serve.Use(cors.New(cors.Options{
		AllowedOrigins:   cfg.Cors.Origin,
		AllowedMethods:   cfg.Cors.Method,
		AllowedHeaders:   cfg.Cors.AllowHeader,
		ExposedHeaders:   cfg.Cors.ExposedHeader,
		MaxAge:           cfg.Cors.MaxAge,
		AllowCredentials: cfg.Cors.Credentials,
	}))
	db, err := bootstrap.InitializeDatabase(&cfg.Mysql)
	if err != nil {
		log.Fatalln(err)
	}
	rdb := bootstrap.InitializeRedis(&cfg.Redis)
	// Define shared variables
	token.Key = []byte(cfg.App.Key)
	token.Options = cfg.Token
	// Configure containers
	serve.ConfigureContainer(func(container *router.APIContainer) {
		container.RegisterDependency(db)
		container.RegisterDependency(rdb)
		container.RegisterDependency(cache.Initialize)
		container.RegisterDependency(curd.Initialize)
		container.Get("/", route.Default)
		container.Options("*", route.Default)
		mvc.Configure(container.Party("/").Self, app.Application)
	})
	serve.Listen(
		cfg.Listen,
		iris.WithConfiguration(iris.YAML("./config/iris.yml")),
	)
}
