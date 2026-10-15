package main

import (
	"fmt"
	"os"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	if err := generate("./config/values.yml",
		func(g *gen.Generator) {
			g.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
				"jsonb": func(columnType gorm.ColumnType) (dataType string) {
					return "common.M"
				},
			})
			g.ApplyBasic(
				g.GenerateModelAs("org", "Org"),
				g.GenerateModelAs("user_org", "UserOrg"),
				g.GenerateModelAs("user_role", "UserRole"),
				g.GenerateModelAs("permission", "Permission"),
				g.GenerateModelAs("role", "Role",
					gen.FieldType("strategy", "common.M"),
				),
				g.GenerateModelAs("route", "Route",
					gen.FieldType("strategy", "common.M"),
				),
				g.GenerateModelAs("user", "User"),
			)
		},
		gen.FieldType("id", "string"),
		gen.FieldType("pid", "string"),
		gen.FieldType("role_id", "string"),
		gen.FieldType("user_id", "string"),
		gen.FieldType("org_id", "string"),
		gen.FieldType("contact_id", "string"),
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
