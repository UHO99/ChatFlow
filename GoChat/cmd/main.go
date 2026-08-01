package main

import (
	"GoChat/api"
	"GoChat/config"
	"log"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load Config : ", err)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("Cannot create server : ", err)
	}

	log.Print(server)
}
