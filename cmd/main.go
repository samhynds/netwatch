package main

import (
	"log"
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/crawl"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// defer cleanup()

	log.Println("Starting...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 1. Load and parse the .netwatch file provided by the user
	cfg, err := config.Load("./sites.netwatch")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Set up channels and Crawl Manager
	queue, _ := crawl.NewManager(cfg)

	// 3. Set up crawl workers
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go func() {
			for url := range queue {
				go crawl.Worker(url, cfg, queue)
			}
		}()
	}

	time.Sleep(time.Second)
	queue <- "meowdy"
	time.Sleep(time.Second)

	<-sigChan

	log.Println("Shutting down...")
	// Perform cleanup then...
	os.Exit(0)

	// // 5. Set up the transport queue
	// transportQueue := make(chan string, 100)

	// // 6. Set up the transport manager
	// transportManager := transport.NewTransportManager(&transportQueue)

	// // 7. Set up transport workers
	// for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
	// 	go transportWorker(transportManager)
	// }
}
