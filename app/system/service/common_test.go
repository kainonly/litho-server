package service

import (
	"github.com/caarlos0/env/v6"
	"lab-api/bootstrap"
	"log"
	"os"
	"testing"
)

var index *Index
var resource *Resource

func TestMain(m *testing.M) {
	d := new(Dependency)
	var err error
	if err = env.Parse(&d.Config); err != nil {
		log.Fatalln(err)
	}
	if d.Db, err = bootstrap.InitializeDatabase(d.Config); err != nil {
		log.Fatalln(err)
	}
	if d.Redis, err = bootstrap.InitializeRedis(d.Config); err != nil {
		log.Fatalln(err)
	}
	index = NewIndex(d)
	resource = NewResource(d)
	os.Exit(m.Run())
}
