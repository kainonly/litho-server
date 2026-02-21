package common

type ActionDef struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

var ActionLabels = map[string]ActionDef{
	"create":      {Label: "新增", Value: "create"},
	"update":      {Label: "更新", Value: "update"},
	"delete":      {Label: "删除", Value: "delete"},
	"sort":        {Label: "排序", Value: "sort"},
	"regroup":     {Label: "重新分组", Value: "regroup"},
	"set_roles":   {Label: "设置角色", Value: "set_roles"},
	"set_actives": {Label: "批量启用/禁用", Value: "set_actives"},
}
