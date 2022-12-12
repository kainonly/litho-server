package model_test

import (
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/common"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var values *common.Values
var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	values, err := bootstrap.LoadStaticValues()
	if err != nil {
		panic(err)
	}

	if db, err = bootstrap.UseDatabase(values); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}
