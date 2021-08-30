package model

import (
	"github.com/caarlos0/env/v6"
	"gorm.io/gorm"
	"lab-api/bootstrap"
	"lab-api/common"
	"log"
	"os"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	var cfg common.Config
	if err = env.Parse(&cfg); err != nil {
		log.Fatalln(err)
	}
	if db, err = bootstrap.InitializeDatabase(cfg); err != nil {
		return
	}
	os.Exit(m.Run())
}
