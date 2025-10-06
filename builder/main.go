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
				g.GenerateModelAs("user", "User"),
			)
		},
		gen.FieldType("id", "string"),
		gen.FieldType("user_id", "string"),
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
