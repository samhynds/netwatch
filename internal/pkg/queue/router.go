package queue

import (
	"log"
	"netwatch/internal/pkg/config"
	"time"
)

/*

	When global limit reached, workers sleep until limit not met.

	- Check when the last request to a host was made
		- If timestamp and count allows for more requests, then make request
		- Otherwise, re-add item to cooldown queue and don't change timestamp or req count

		To do that, we need to save per-host timestamps and req count alongside the queue.

*/

// Returns the next item to crawl and the name of the queue that item came from
func (q *Queues) GetNextItem(config *config.Config) (string, any) {
	// if item in recrawl, return that first if time is correct
	if len(q.Recrawl.Queue) > 0 {
		for item := range q.Recrawl.Get() {
			if item.TimeCrawlable.Before(time.Now()) {
				return item.URL, &q.Recrawl
			} else {
				q.Recrawl.Queue <- item
			}
		}
	}

	// if items in cooldown, check and return first one we can crawl (time is correct)
	if len(q.Cooldown.Queue) > 0 {
		for item := range q.Cooldown.Get() {
			if item.TimeCrawlable.Before(time.Now()) {
				host, err := GetHostFromURL(item.URL)
				if err != nil {
					log.Println("Error getting host for", item.URL)
					continue
				}

				q.processedHosts[host] = ProcessedHost{
					count: q.processedHosts[host].count - 1,
				}

				return item.URL, &q.Cooldown
			} else {
				q.Cooldown.Queue <- item
			}
		}
	}

	// if recrawl or cooldown don't return an item to crawl, return the next one from crawl queue
	return <-q.Crawl.Get(), &q.Crawl
}
