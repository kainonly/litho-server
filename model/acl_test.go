package model

import "testing"

func TestAcl(t *testing.T) {
	if err := db.AutoMigrate(&Acl{}); err != nil {
		t.Error(err)
	}
	data := []Acl{
		{
			Name: "访问控制",
			Path: "acl",
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
			Name: "资源控制",
			Path: "resource",
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
			Name: "特殊授权",
			Path: "permission",
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
			Name: "权限组",
			Path: "role",
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
			Name: "成员管理",
			Path: "admin",
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
