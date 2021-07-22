package main

import (
	"github.com/kainonly/gin-helper/hash"
	"lab-api/bootstrap"
	"lab-api/model"
	"log"
)

func main() {
	cfg, err := bootstrap.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := bootstrap.InitializeDatabase(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	db = db.Debug()
	db.AutoMigrate(&model.Admin{})
	password, _ := password.Make("pass@VAN1234")
	db.Create(&model.Admin{
		Username: "kain",
		Password: password,
		Super:    model.True(),
	})
}
