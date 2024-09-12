package crawl

import "netwatch/internal/pkg/config"

func NewManager(config *config.Config) (*CrawlQueue, chan string) {
	// Sites that are ready to be crawled
	var queue = NewCrawlQueue(10)

	// Sites to be recrawled periodically if enabled
	var recrawlQueue = make(chan string, 100)

	// Load the initial sites from the config file into the queue
	for _, site := range config.Sites {
		queue.Add(site.URL)

		if config.Config.Recrawl.Enabled {
			recrawlQueue <- site.URL
		}
	}

	return queue, recrawlQueue
}
