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

func Worker(
	url string,
	loadedConfig *config.Config,
	queue *q.Queue[string],
	rateLimiter *ratelimiter.RateLimiter,
	i int,
) (q.ProcessedQueueItem, error) {
	defer queue.MarkProcessed(url)
	log.Printf("Worker #%d started with URL: %s", i, url)
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
		return q.ProcessedQueueItem{}, err
	}

	log.Println("Site config for", url, "is", siteConfig)

	// Make Request
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error making request to", url)
		log.Println(err)
		return q.ProcessedQueueItem{}, err
	}

	defer res.Body.Close()

	// Read Response Body
	var bodyBuffer bytes.Buffer
	io.ReadAll(io.TeeReader(res.Body, &bodyBuffer))

	// Extract Links & Content
	doc, err := NewContentExtractor(io.NopCloser(&bodyBuffer))
	if err != nil {
		log.Println("Error creating content extractor for", url)
		return q.ProcessedQueueItem{}, err
	}

	links, err := LinkExtractor(doc, url, loadedConfig.Config.Roam, &siteConfig.Links)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return q.ProcessedQueueItem{}, err
	}

	content, err := ContentExtractor(doc, &siteConfig.Content)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return q.ProcessedQueueItem{}, err
	}

	docHtml, err := doc.Html()
	if err != nil {
		log.Println("Error parsing response body for", url)
		return q.ProcessedQueueItem{}, err
	}

	// fmt.Printf("\nLinks for %s: %s", url, links)
	// fmt.Printf("\nContent for %s: %s", url, content)

	return q.ProcessedQueueItem{
		URL:       url,
		Timestamp: time.Now(),
		Content:   content,
		Links:     links,
		Body:      docHtml,
		Headers:   res.Header,
	}, nil
}
