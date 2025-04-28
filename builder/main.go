package main

import (
	"fmt"
	"os"
)

func main() {
	if err := generate("./config/values.yml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
