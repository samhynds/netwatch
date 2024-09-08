package main

import (
	"log"
	"netwatch/internal/pkg/config"
)

func main() {
	log.Println("NetWatch is starting...")

	// 1. Load the .netwatch file provided by the user
	cfg, err := config.Load("./sites.netwatch")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	// 2. Parse the .netwatch file

}
