package crawl

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"netwatch/internal/pkg/config"
	"time"
)

func Worker(url string, loadedConfig *config.Config, queue *CrawlQueue) ([]byte, error) {
	defer queue.MarkProcessed(url)
	log.Println("Worker started with URL: ", url)
	time.Sleep(time.Second)

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
	body, err := io.ReadAll(io.TeeReader(res.Body, &bodyBuffer))
	if err != nil {
		log.Println("Error reading response body for", url)
		return nil, err
	}

	// Extract Links
	links, err := LinkExtractor(io.NopCloser(&bodyBuffer), loadedConfig.Config.Roam, &siteConfig.Links)
	if err != nil {
		log.Println("Error parsing response body for", url)
		return nil, err
	}

	// Send links to crawl queue - add link dedupe later
	for _, link := range links {
		log.Println("Adding:", link)
		queue.Add(link)
	}

	// TODO: Extract content from site

	fmt.Printf("\nLinks for %s: %s", url, links)
	fmt.Printf("\nBody for %s: %s", url, string(body[0:50]))

	// return or chan?
	return nil, nil
}
