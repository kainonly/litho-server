package role_resource_rel

import (
	"gorm.io/gorm"
	"taste-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.RoleResourceRel{}); err != nil {
		return err
	}
	data := []model.RoleResourceRel{
		{RoleKey: "*", ResourceKey: "system"},
		{RoleKey: "*", ResourceKey: "center"},
		{RoleKey: "*", ResourceKey: "profile"},
		{RoleKey: "*", ResourceKey: "acl-index"},
		{RoleKey: "*", ResourceKey: "acl-add"},
		{RoleKey: "*", ResourceKey: "acl-edit"},
		{RoleKey: "*", ResourceKey: "admin-index"},
		{RoleKey: "*", ResourceKey: "admin-add"},
		{RoleKey: "*", ResourceKey: "admin-edit"},
		{RoleKey: "*", ResourceKey: "resource-index"},
		{RoleKey: "*", ResourceKey: "resource-add"},
		{RoleKey: "*", ResourceKey: "resource-edit"},
		{RoleKey: "*", ResourceKey: "role-index"},
		{RoleKey: "*", ResourceKey: "role-add"},
		{RoleKey: "*", ResourceKey: "role-edit"},
	}
	return db.Create(&data).Error
}
