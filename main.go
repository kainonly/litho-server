package main

import (
	"lab-api/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	a, err := Bootstrap(app)
	if err != nil {
		log.Fatalln(err)
	}

	a.Run(":9000")
}
