package main

import (
	"github.com/caarlos0/env/v6"
	"lab-api/common"
	"log"
)

func main() {
	var cfg common.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln(err)
	}
	app, err := App(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	app.Run(":9000")
}
