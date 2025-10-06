package main

import (
	"context"
	"fmt"
	"os"
	"server/model"

	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passlib"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := generate("./config/values.yml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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
	db.Exec("SET client_encoding TO 'UTF8'")

	ctx := context.Background()
	data := model.User{
		ID:    help.SID(),
		Email: "work@kainonly.com",
	}
	data.Password, _ = passlib.Hash(`pass@VAN1234`)

	if err = db.WithContext(ctx).Create(&data).Error; err != nil {
		return
	}
	return
}
