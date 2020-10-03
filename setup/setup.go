package main

import (
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"os"
	"van-api/app/model"
	"van-api/app/types"
)

func main() {
	if _, err := os.Stat("./config/config.yml"); os.IsNotExist(err) {
		log.Fatalln("the configuration file does not exist")
	}
	buf, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("failed to read service configuration file", err)
	}
	config := types.Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatalln("service configuration file parsing failed", err)
	}
	db, err := gorm.Open(mysql.Open(config.Mysql.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Mysql.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Api Access Control Table'",
	).AutoMigrate(&model.Acl{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Resource Access Control Table'",
	).AutoMigrate(&model.Resource{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Policy Table'",
	).AutoMigrate(&model.Policy{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Role Table'",
	).AutoMigrate(&model.RoleBasic{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Role Associated Resource Table'",
	).AutoMigrate(&model.RoleResourceAssoc{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Admin Table'",
	).AutoMigrate(&model.AdminBasic{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Set(
		"gorm:table_options",
		"comment='Admin Associated Role Table'",
	).AutoMigrate(&model.AdminRoleAssoc{})
	if err != nil {
		log.Fatalln(err)
	}
}
