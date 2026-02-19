package main

import (
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Values struct {
	Database DatabaseValue `yaml:"database"`
}

type DatabaseValue struct {
	Debug bool   `yaml:"debug"`
	DSN   string `yaml:"dsn"`
	Name  string `yaml:"name"`
}

func generate(path string, fn func(g *gen.Generator), opts ...gen.ModelOpt) (err error) {
	os.RemoveAll("./model")
	var v *Values
	v = new(Values)
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(b, &v); err != nil {
		return
	}
	var db *gorm.DB
	if db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  v.Database.DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{}); err != nil {
		return
	}
	db.Exec("SET client_encoding TO 'UTF8'")
	g := gen.NewGenerator(gen.Config{
		OutPath:          "./query",
		ModelPkgPath:     "./model",
		Mode:             gen.WithDefaultQuery,
		FieldNullable:    true,
		FieldCoverable:   true,
		FieldSignable:    true,
		FieldWithTypeTag: true,
	})

	g.WithOpts(opts...)
	g.UseDB(db)
	fn(g)
	g.Execute()
	os.RemoveAll("./query")
	return
}
