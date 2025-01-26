package main

import (
	"flag"
	"log"
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/crawl"
	"netwatch/internal/pkg/queue"
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

	// Load and parse the .netwatch file provided by the user
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Set up channels for crawling and transporting
	queues := queue.NewQueue(cfg)
	queues.InitPopulation(cfg)

	// Set up rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(
		cfg.Config.Requests.MaxTotal,
		cfg.Config.Requests.Window,
	)

	// Start crawl workers
	log.Println("Starting", cfg.Config.Requests.MaxConcurrent, "crawl workers")
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go func(workerID int) {
			for {
				url, queueSource := queues.GetNextItem(cfg)

				formattedResponse, err := crawl.Worker(url, cfg, queueSource, &queues.Crawl, rateLimiter, workerID)

				if err != nil {
					log.Printf("Worker %d error: %v\n", workerID, err)
					continue
				}

				queues.Crawl.AddMultiple(formattedResponse.Links)
				queues.Transport.Add(formattedResponse)
			}
		}(i)
	}

	// Set up transporter connections
	dbConn, kafkaConn := transporter.SetupConnections(cfg)

	// Start transport workers
	for i := 0; i < cfg.Config.Requests.MaxConcurrent; i++ {
		go func() {
			for item := range queues.Transport.Get() {
				log.Println("Transporting", item.URL)
				transporter.Worker(item, dbConn, kafkaConn)
			}
		}()
	}

	// Shutdown gracefully on SIGINT or SIGTERM
	<-sigChan
	log.Println("Shutting down...")
	os.Exit(0)
}
