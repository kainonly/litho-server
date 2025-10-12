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
				g.GenerateModelAs("job", "Job"),
				g.GenerateModelAs("scheduler", "Scheduler",
					gen.FieldType("thumbs", "common.M"),
				),
				g.GenerateModelAs("team", "Team"),
				g.GenerateModelAs("team_user", "TeamUser"),
				g.GenerateModelAs("user", "User"),
			)
		},
		gen.FieldType("id", "string"),
		gen.FieldType("user_id", "string"),
		gen.FieldType("team_id", "string"),
		gen.FieldType("scheduler_id", "string"),
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
