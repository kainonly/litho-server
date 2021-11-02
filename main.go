package main

import "log"

func main() {
	serve, err := API()
	if err != nil {
		log.Panicln(err)
	}
	serve.Run(":9000")
}
