package main

import (
	"lab-api/bootstrap"
	"log"
)

func main() {
	cfg, err := bootstrap.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	app, err := App(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	app.Run(":9000")
}
