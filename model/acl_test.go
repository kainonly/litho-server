package model

import "testing"

func TestAcl(t *testing.T) {
	if err := db.AutoMigrate(&Acl{}); err != nil {
		t.Error(err)
	}
	data := []Acl{
		{
			Key:  "acl",
			Name: "访问控制",
			Acts: Acts{
				R: Act{
					"originLists": "获取原始列表资源",
					"lists":       "获取分页列表资源",
					"get":         "获取单条资源",
				},
				W: Act{
					"add":    "创建资源",
					"edit":   "更新资源",
					"delete": "删除资源",
				},
			},
		},
		{
			Key:  "resource",
			Name: "资源控制",
			Acts: Acts{
				Act{
					"originLists": "获取原始列表资源",
					"get":         "获取单条资源",
				},
				Act{
					"add":    "创建资源",
					"edit":   "更新资源",
					"delete": "删除资源",
				},
			},
		},
		{
			Key:  "role",
			Name: "权限组",
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
		{
			Key:  "admin",
			Name: "用户管理",
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
