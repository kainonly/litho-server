package main

import (
	"github.com/weplanx/server/bootstrap"
	"github.com/weplanx/server/utils/validation"
)

func main() {
	values, err := bootstrap.LoadStaticValues("./config/config.yml")
	if err != nil {
		panic(err)
	}

	validation.Extend()

	server, err := bootstrap.OkLetsGo(values)
	if err != nil {
		panic(err)
	}

	server.Spin()
}
