package resource

import (
	"gorm.io/gorm"
	"taste-api/application/common/types"
	"taste-api/application/model"
)

func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Resource{}); err != nil {
		return err
	}
	data := []model.Resource{
		{Key: "center", Parent: "origin", Name: types.JSON{"zh_cn": "个人中心", "en_us": "Center"}},
		{Key: "profile", Parent: "center", Name: types.JSON{"zh_cn": "信息修改", "en_us": "Profile"}, Router: true},
		{Key: "system", Parent: "origin", Name: types.JSON{"zh_cn": "系统设置", "en_us": "System"}, Nav: true, Icon: "setting"},
		{Key: "resource-index", Parent: "system", Name: types.JSON{"zh_cn": "资源控制管理", "en_us": "Resource"}, Nav: true, Router: true, Policy: true},
		{Key: "resource-add", Parent: "resource-index", Name: types.JSON{"zh_cn": "资源控制新增", "en_us": "Resource Add"}, Router: true},
		{Key: "resource-edit", Parent: "resource-index", Name: types.JSON{"zh_cn": "资源控制修改", "en_us": "Resource Edit"}, Router: true},
		{Key: "acl-index", Parent: "system", Name: types.JSON{"zh_cn": "访问控制管理", "en_us": "Acl"}, Nav: true, Router: true, Policy: true},
		{Key: "acl-add", Parent: "acl-index", Name: types.JSON{"zh_cn": "访问控制新增", "en_us": "Acl Add"}, Router: true},
		{Key: "acl-edit", Parent: "acl-index", Name: types.JSON{"zh_cn": "访问控制修改", "en_us": "Acl Edit"}, Router: true},
		{Key: "role-index", Parent: "system", Name: types.JSON{"zh_cn": "权限组", "en_us": "Role"}, Nav: true, Router: true, Policy: true},
		{Key: "role-add", Parent: "role-index", Name: types.JSON{"zh_cn": "权限组新增", "en_us": "Role Add"}, Router: true},
		{Key: "role-edit", Parent: "role-index", Name: types.JSON{"zh_cn": "权限组修改", "en_us": "Role Edit"}, Router: true},
		{Key: "admin-index", Parent: "system", Name: types.JSON{"zh_cn": "管理员", "en_us": "Admin"}, Nav: true, Router: true, Policy: true},
		{Key: "admin-add", Parent: "admin-index", Name: types.JSON{"zh_cn": "管理员新增", "en_us": "Admin Add"}, Router: true},
		{Key: "admin-edit", Parent: "admin-index", Name: types.JSON{"zh_cn": "管理员修改", "en_us": "Admin Edit"}, Router: true},
	}
	return db.Create(&data).Error
}
