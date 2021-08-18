package model

import "testing"

func TestAcl(t *testing.T) {
	if err := db.AutoMigrate(&Acl{}); err != nil {
		t.Error(err)
	}
	data := []Acl{
		{
			Key:  "acl",
			Name: "访问控制权",
			Acts: [2][]Act{
				{
					{Path: "originLists", Description: "获取原始列表资源"},
					{Path: "lists", Description: "获取分页列表资源"},
					{Path: "get", Description: "获取单条资源"},
				},
				{
					{Path: "add", Description: "创建资源"},
					{Path: "edit", Description: "更新资源"},
					{Path: "delete", Description: "删除资源"},
				},
			},
		},
		{
			Key:  "resource",
			Name: "资源控制权",
			Acts: [2][]Act{
				{
					{Path: "originLists", Description: "获取原始列表资源"},
					{Path: "get", Description: "获取单条资源"},
				},
				{
					{Path: "add", Description: "创建资源"},
					{Path: "edit", Description: "更新资源"},
					{Path: "delete", Description: "删除资源"},
				},
			},
		},
	}
	if err := db.Create(&data).Error; err != nil {
		t.Error(err)
	}

}
