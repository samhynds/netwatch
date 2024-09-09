package crawl

import "netwatch/internal/pkg/config"

func NewManager(config *config.Config) (chan string, chan string) {
	// Sites that are ready to be crawled
	var queue = make(chan string, 100000) // (1.6mb empty)

	// Sites to be recrawled periodically if enabled
	var recrawlQueue = make(chan string, 100)

	// Load the initial sites from the config file into the queue
	for _, site := range config.Sites {
		queue <- site.URL

		if config.Config.Recrawl.Enabled {
			recrawlQueue <- site.URL
		}
	}

	return queue, recrawlQueue
}

func Manager() {
	// return the next item in either the queue or recrawl queue

}
