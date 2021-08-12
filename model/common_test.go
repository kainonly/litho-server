package model

import (
	"github.com/kainonly/go-bit"
	"gorm.io/gorm"
	"lab-api/bootstrap"
	"log"
	"os"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	os.Chdir(`../`)
	config, err := bit.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	if db, err = bootstrap.InitializeDatabase(config); err != nil {
		return
	}
	db = db.Debug()
	os.Exit(m.Run())
}
