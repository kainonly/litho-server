package acl

import (
	"gorm.io/gorm"
	"taste-api/application/common/types"
	"taste-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Acl{}); err != nil {
		return err
	}
	data := []model.Acl{
		{Key: "main", Name: types.JSON{"zh_cn": "公共模块", "en_us": "Common Module"}, Write: "uploads", Read: ""},
		{Key: "resource", Name: types.JSON{"zh_cn": "资源控制模块", "en_us": "Resource Module"}, Write: "add,edit,delete,sort", Read: "originLists,lists,get"},
		{Key: "acl", Name: types.JSON{"zh_cn": "访问控制模块", "en_us": "Acl Module"}, Write: "add,edit,delete", Read: "originLists,lists,get"},
		{Key: "policy", Name: types.JSON{"zh_cn": "策略模块", "en_us": "Policy Module"}, Write: "add,delete", Read: "originLists"},
		{Key: "admin", Name: types.JSON{"zh_cn": "管理员模块", "en_us": "Admin Module"}, Write: "add,edit,delete", Read: "originLists,lists,get"},
		{Key: "role", Name: types.JSON{"zh_cn": "权限组模块", "en_us": "Role Module"}, Write: "add,edit,delete", Read: "originLists,lists,get"},
	}
	return db.Create(&data).Error
}
