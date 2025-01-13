package crawl

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/ratelimiter"
	"time"
)

func Worker(url string, loadedConfig *config.Config, queue *CrawlQueue, rateLimiter *ratelimiter.RateLimiter, i int) ([]byte, error) {
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
		return nil, err
	}

	log.Println("Site config for", url, "is", siteConfig)

	// Make Request
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error making request to", url)
		log.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	// Read Response Body
	var bodyBuffer bytes.Buffer
	io.ReadAll(io.TeeReader(res.Body, &bodyBuffer))

	// Extract Links & Content
	doc, err := NewContentExtractor(io.NopCloser(&bodyBuffer))
	if err != nil {
		log.Println("Error creating content extractor for", url)
		return nil, err
	}

	links, err := LinkExtractor(doc, url, loadedConfig.Config.Roam, &siteConfig.Links)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return nil, err
	}

	queue.AddMultiple(links)

	content, err := ContentExtractor(doc, &siteConfig.Content)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return nil, err
	}

	fmt.Printf("\nLinks for %s: %s", url, links)
	fmt.Printf("\nContent for %s: %s", url, content)

	// return or chan?
	return nil, nil
}
