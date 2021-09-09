package main

import (
	"lab-api/common"
	"log"
)

func main() {
	set, err := common.LoadSettings()
	if err != nil {
		log.Fatalln(err)
	}
	app, err := App(set)
	if err != nil {
		log.Fatalln(err)
	}
	app.Run(":9000")
}
