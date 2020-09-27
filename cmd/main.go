package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"van-api/app"
	"van-api/app/types"
)

func main() {
	if _, err := os.Stat("./config/config.yml"); os.IsNotExist(err) {
		log.Fatalln("the configuration file does not exist")
	}
	buf, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalln("failed to read service configuration file", err)
	}
	config := types.Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatalln("service configuration file parsing failed", err)
	}
	log.Println(config)
	app.Application(&config)
}
