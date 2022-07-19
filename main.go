package main

import (
	"github.com/weplanx/server/common"
)

func main() {
	values, err := common.LoadStaticValues("./config/config.yml")
	if err != nil {
		panic(err)
	}

	server, err := OkLetsGo(values)
	if err != nil {
		panic(err)
	}

	server.Run(":3000")
}
