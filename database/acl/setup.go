package acl

import (
	"github.com/kainonly/gin-extra/datatype"
	"gorm.io/gorm"
	"lab-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Acl{}); err != nil {
		return err
	}
	data := []model.Acl{
		{
			Key:   "main",
			Name:  datatype.JSONObject{"zh_cn": "公共模块", "en_us": "Common Module"},
			Write: datatype.JSONArray{"uploads"},
		},
		{
			Key:   "resource",
			Name:  datatype.JSONObject{"zh_cn": "资源控制模块", "en_us": "Resource Module"},
			Write: datatype.JSONArray{"add", "edit", "delete", "sort"},
			Read:  datatype.JSONArray{"originLists", "lists", "get"},
		},
		{
			Key:   "acl",
			Name:  datatype.JSONObject{"zh_cn": "访问控制模块", "en_us": "Acl Module"},
			Write: datatype.JSONArray{"add", "edit", "delete"},
			Read:  datatype.JSONArray{"originLists", "lists", "get"},
		},
		{
			Key:   "policy",
			Name:  datatype.JSONObject{"zh_cn": "策略模块", "en_us": "Policy Module"},
			Write: datatype.JSONArray{"add", "delete"},
			Read:  datatype.JSONArray{"originLists"},
		},
		{
			Key:   "admin",
			Name:  datatype.JSONObject{"zh_cn": "管理员模块", "en_us": "Admin Module"},
			Write: datatype.JSONArray{"add", "edit", "delete"},
			Read:  datatype.JSONArray{"originLists", "lists", "get"},
		},
		{
			Key:   "role",
			Name:  datatype.JSONObject{"zh_cn": "权限组模块", "en_us": "Role Module"},
			Write: datatype.JSONArray{"add", "edit", "delete"},
			Read:  datatype.JSONArray{"originLists", "lists", "get"},
		},
	}
	return db.Create(&data).Error
}
