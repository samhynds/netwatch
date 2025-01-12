package main

import (
	"flag"
	"log"
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/crawl"
	"netwatch/internal/pkg/ratelimiter"
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

	configFile := flag.String("config", "./sites.netwatch", "Path to the .netwatch configuration file")
	flag.Parse()

	// 1. Load and parse the .netwatch file provided by the user
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Set up channels and Crawl Manager
	manager := crawl.NewManager(cfg)

	// 3. Set up rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(
		cfg.Config.Requests.MaxTotal,
		cfg.Config.Requests.Window,
	)

	log.Println("Starting", cfg.Config.Requests.MaxConcurrent, "crawl workers...")
	time.Sleep(time.Second)
	// 3. Set up crawl workers
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go func() {
			for url := range manager.Queue.Get() {
				crawl.Worker(url, cfg, manager.Queue, rateLimiter, i)
			}
		}()
	}

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
