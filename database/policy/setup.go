package policy

import (
	"gorm.io/gorm"
	"lab-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Policy{}); err != nil {
		return err
	}
	data := []model.Policy{
		{ResourceKey: "acl-index", AclKey: "acl", Policy: 1},
		{ResourceKey: "resource-index", AclKey: "resource", Policy: 1},
		{ResourceKey: "resource-index", AclKey: "policy", Policy: 1},
		{ResourceKey: "resource-index", AclKey: "acl", Policy: 0},
		{ResourceKey: "role-index", AclKey: "role", Policy: 1},
		{ResourceKey: "role-index", AclKey: "resource", Policy: 0},
		{ResourceKey: "admin-index", AclKey: "admin", Policy: 1},
		{ResourceKey: "admin-index", AclKey: "role", Policy: 0},
	}
	return db.Create(&data).Error
}
