package main

import (
	"api/bootstrap"
	"github.com/gin-contrib/pprof"
	"os"
)

func main() {
	values, err := bootstrap.SetValues()
	if err != nil {
		panic(err)
	}
	app, err := App(values)
	if err != nil {
		panic(err)
	}
	if os.Getenv("GIN_MODE") != "release" {
		pprof.Register(app)
	}
	app.Run(values.Address)
}
