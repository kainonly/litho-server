package role_basic

import (
	"gorm.io/gorm"
	"lab-api/application/common/datatype"
	"lab-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.RoleBasic{}); err != nil {
		return err
	}
	data := model.RoleBasic{
		Key:  "*",
		Name: datatype.JSONObject{"zh_cn": "超级管理员", "en_us": "super"},
	}
	return db.Create(&data).Error
}
