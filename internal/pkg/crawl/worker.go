package crawl

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"netwatch/internal/pkg/config"
	q "netwatch/internal/pkg/queue"
	"netwatch/internal/pkg/ratelimiter"
	"time"
)

func Worker[T string | q.CrawledQueueItem | q.ScheduledQueueItem](
	url string,
	loadedConfig *config.Config,
	queues *q.Queues,
	queueForItem *q.Queue[T],
	rateLimiter *ratelimiter.RateLimiter,
	i int,
) (q.CrawledQueueItem, error) {
	defer queueForItem.MarkProcessed(url)
	log.Printf("Worker #%d started with URL: %s", i, url)

	// Check if per host limit reached
	// Check if queues.ProcessedHosts[host].count < config.maxPerHost
	// If not, add url to queues.Cooldown

	// Check rate limit, sleep if needed
	allowed, nextAllowedTime := rateLimiter.Allow()
	log.Printf("Allowed: %t, Next allowed time: %v", allowed, nextAllowedTime)
	if !allowed {
		waitDuration := time.Until(nextAllowedTime)
		log.Printf("Rate limit exceeded for %s. Next allowed time: %v (in %v)", url, nextAllowedTime, waitDuration)

		// TODO: This consumes a worker for waitDuration, instead schedule this item to be processed later and continue with other items now
		time.Sleep(waitDuration)
	}

	// Find config that matches the url as a regex
	siteConfig, err := config.SiteConfigForURL(url, loadedConfig)
	if err != nil {
		log.Println("Error finding site config for", url)
		return q.CrawledQueueItem{}, err
	}

	log.Println("Site config for", url, "is", siteConfig)

	// Make Request
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error making request to", url)
		log.Println(err)
		return q.CrawledQueueItem{}, err
	}

	defer res.Body.Close()

	// Read Response Body
	var bodyBuffer bytes.Buffer
	io.ReadAll(io.TeeReader(res.Body, &bodyBuffer))

	// Extract Links & Content
	doc, err := NewContentExtractor(io.NopCloser(&bodyBuffer))
	if err != nil {
		log.Println("Error creating content extractor for", url)
		return q.CrawledQueueItem{}, err
	}

	links, err := LinkExtractor(doc, url, loadedConfig.Config.Roam, &siteConfig.Links)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return q.CrawledQueueItem{}, err
	}

	content, err := ContentExtractor(doc, &siteConfig.Content)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return q.CrawledQueueItem{}, err
	}

	docHtml, err := doc.Html()
	if err != nil {
		log.Println("Error parsing response body for", url)
		return q.CrawledQueueItem{}, err
	}

	// fmt.Printf("\nLinks for %s: %s", url, links)
	// fmt.Printf("\nContent for %s: %s", url, content)

	return q.CrawledQueueItem{
		URL:       url,
		Timestamp: time.Now(),
		Content:   content,
		Links:     links,
		Body:      docHtml,
		Headers:   res.Header,
	}, nil
}
