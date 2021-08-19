package model

import "testing"

func TestPermission(t *testing.T) {
	if err := db.AutoMigrate(&Permission{}); err != nil {
		t.Error(err)
	}
	data := []Permission{
		{
			Code: "ACCESS_ACL",
			Name: "访问控制",
		},
		{
			Code: "ACCESS_RESOURCE",
			Name: "资源控制",
		},
		{
			Code: "ACCESS_ROLE",
			Name: "权限组",
		},
		{
			Code: "ACCESS_ADMIN",
			Name: "用户管理",
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
