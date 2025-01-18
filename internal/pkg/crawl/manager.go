package crawl

import "netwatch/internal/pkg/config"

type Manager struct {
	Queue        *CrawlQueue
	RecrawlQueue chan string
}

func NewManager(config *config.Config) *Manager {
	// Sites that are ready to be crawled (1m = ~16mb memory empty)
	queue := NewCrawlQueue(config.Config.Queue.Capacity)

	// Sites to be recrawled periodically if enabled
	recrawlQueue := make(chan string, 100)

	// Load the initial sites from the config file into the queue
	for _, site := range config.Sites {
		queue.Add(site.URL)

		if config.Config.Recrawl.Enabled {
			recrawlQueue <- site.URL
		}
	}

	return &Manager{
		Queue:        queue,
		RecrawlQueue: recrawlQueue,
	}
}
