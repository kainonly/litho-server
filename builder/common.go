package main

import (
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"server/model"
)

type Values struct {
	Database `yaml:"database"`
}

type Database struct {
	Url string `yaml:"url"`
}

func generate(path string) (err error) {
	var v *Values
	v = new(Values)
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &v); err != nil {
		return
	}
	var db *gorm.DB
	if db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  v.Database.Url,
		PreferSimpleProtocol: true,
	}), &gorm.Config{}); err != nil {
		return
	}
	db.AutoMigrate(
		&model.Category{},
		&model.CategoryRef{},
		&model.Department{},
		&model.Picture{},
		&model.Resource{},
		&model.Role{},
		&model.Route{},
		&model.User{},
		&model.Video{},
	)
	return
}
