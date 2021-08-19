package model

import "testing"

func TestAcl(t *testing.T) {
	if err := db.AutoMigrate(&Acl{}); err != nil {
		t.Error(err)
	}
	data := []Acl{
		{
			Name:  "用户管理",
			Model: "admin",
			Acts: Acts{
				Act{
					"originLists": "获取原始列表资源",
					"lists":       "获取分页列表资源",
					"get":         "获取单条资源",
				},
				Act{
					"add":    "创建资源",
					"edit":   "更新资源",
					"delete": "删除资源",
				},
			},
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}
}
