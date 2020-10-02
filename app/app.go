package app

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"van-api/app/controller"
	"van-api/app/middleware/cors"
	"van-api/app/types"
)

func Application(option *types.Config) {
	app := iris.Default()
	app.Use(cors.Cors(types.CorsOption{}))
	context.DefaultJSONOptions = context.JSON{
		StreamingJSON: true,
	}
	db, err := gorm.Open(mysql.Open(option.Mysql.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	if option.Mysql.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(option.Mysql.MaxIdleConns)
	}
	if option.Mysql.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(option.Mysql.MaxOpenConns)
	}
	if option.Mysql.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(option.Mysql.ConnMaxLifetime))
	}
	control := controller.New()
	app.Get("/", control.Default)
	app.Options("*", control.Option)
	main := app.Party("/main")
	{
		main.Post("/verify", control.MainVerify)
	}
	app.Listen(option.Listen)
}
