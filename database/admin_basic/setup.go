package admin_basic

import (
	"github.com/kainonly/gin-extra/hash"
	"gorm.io/gorm"
	"lab-api/application/model"
)

func Setup(db *gorm.DB) (err error) {
	if err = db.AutoMigrate(&model.AdminBasic{}); err != nil {
		return err
	}
	var password string
	if password, err = hash.Make("pass@VAN1234"); err != nil {
		return
	}
	return db.Create(&model.AdminBasic{
		Username: "kain",
		Password: password,
	}).Error
}
