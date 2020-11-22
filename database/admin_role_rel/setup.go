package admin_role_rel

import (
	"gorm.io/gorm"
	"taste-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.AdminRoleRel{}); err != nil {
		return err
	}
	return db.Create(&model.AdminRoleRel{
		Username: "kain",
		RoleKey:  "*",
	}).Error
}
