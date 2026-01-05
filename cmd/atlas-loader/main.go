package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"server/model"
)

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
}
