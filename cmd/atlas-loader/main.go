package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"server/model"
)

// 表注释映射（表名 -> 注释）
var tableComments = map[string]string{
	"user":                  "用户表",
	"org":                   "组织表",
	"role":                  "权限表",
	"menu":                  "导航表",
	"route":                 "路由表",
	"resource":              "资源表",
	"permission":            "特定授权表",
	"resource_action":       "资源操作表",
	"role_permission":       "权限特定授权表",
	"role_menu":             "权限导航表",
	"role_route":            "权限路由表",
	"user_org_role":         "用户组织权限表",
	"route_resource_action": "路由资源表",
}

func main() {
	// 定义所有需要迁移的模型
	models := []any{
		&model.User{},
		&model.Org{},
		&model.Role{},
		&model.Menu{},
		&model.Route{},
		&model.Resource{},
		&model.Permission{},
		&model.ResourceAction{},
		&model.RolePermission{},
		&model.RoleMenu{},
		&model.RoleRoute{},
		&model.UserOrgRole{},
		&model.RouteResourceAction{},
	}

	// 使用 PostgreSQL 方言生成 schema
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	io.WriteString(os.Stdout, stmts)

	// 追加表注释 SQL
	for table, comment := range tableComments {
		fmt.Fprintf(os.Stdout, "COMMENT ON TABLE \"%s\" IS '%s';\n", table, comment)
	}
}
