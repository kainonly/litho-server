package main

import (
	"api/bootstrap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"log"
	"os"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			bootstrap.SetValues,
			bootstrap.UseDatabase,
			MockPages,
		),
		fx.Invoke(func(result *mongo.InsertManyResult) {
			log.Println(result)
			os.Exit(0)
		}),
	)
}
