package main

import (
	"lab-api/bootstrap"
	"log"
)

func main() {
	set, err := bootstrap.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}
	app, err := Bootstrap(set)
	if err != nil {
		log.Fatalln(err)
	}

	app.Run(":9000")
}
