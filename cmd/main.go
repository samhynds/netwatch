package main

import (
	"flag"
	"log"
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/crawl"
	"netwatch/internal/pkg/ratelimiter"
	"netwatch/internal/pkg/transporter"
	"os"
	"os/signal"
	"syscall"
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

	// 2. Set up channels for crawling and transporting
	crawlManager := crawl.NewManager(cfg)
	transportManager := transporter.NewManager(cfg)

	// 3. Set up rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(
		cfg.Config.Requests.MaxTotal,
		cfg.Config.Requests.Window,
	)

	log.Println("Starting", cfg.Config.Requests.MaxConcurrent, "crawl workers...")
	// 3. Set up crawl workers
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go func() {
			for url := range crawlManager.Queue.Get() {
				formattedResponse, err := crawl.Worker(url, cfg, crawlManager.Queue, rateLimiter, i)
				if err != nil {
					log.Println(err)
					continue
				}

				crawlManager.Queue.AddMultiple(formattedResponse.Links)
				transportManager.Queue.Add(formattedResponse)
			}
		}()
	}

	// // 7. Set up transport workers
	// for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
	// 	go transportWorker(transportManager)
	// }

	<-sigChan

	log.Println("Shutting down...")
	// Perform cleanup then...
	os.Exit(0)
}
