package main

import (
	"log"
	"netwatch/internal/pkg/config"
)

func main() {
	defer cleanup()
	log.Println("NetWatch is starting...")

	// 1. Load and parse the .netwatch file provided by the user
	cfg, err := config.Load("./sites.netwatch")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	// 2. Set up channels for the crawl queue
	crawlQueue := make(chan string, 100)

	// 3. Set up the Crawl Manager
	crawlManager := crawl.NewCrawlManager(&crawlQueue)

	// 4. Set up crawl workers
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go crawlWorker(crawlManager.readyQueue)
	}

	// 5. Set up the transport queue
	transportQueue := make(chan string, 100)

	// 6. Set up the transport manager
	transportManager := transport.NewTransportManager(&transportQueue)

	// 7. Set up transport workers
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go transportWorker(transportManager.readyQueue)
	}
}
