package crawl

import (
	"log"
	"sync"
)

type CrawlQueue struct {
	queue     chan string
	inQueue   map[string]bool
	processed map[string]bool
	mu        sync.RWMutex
}

func NewCrawlQueue(size int) *CrawlQueue {
	return &CrawlQueue{
		queue:     make(chan string, size),
		inQueue:   make(map[string]bool),
		processed: make(map[string]bool),
	}
}

func (cq *CrawlQueue) Add(url string) bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	if cq.inQueue[url] || cq.processed[url] {
		log.Println("Already in queue or processed", url)
		return false
	}

	cq.inQueue[url] = true
	cq.queue <- url
	return true
}

func (cq *CrawlQueue) Get() <-chan string {
	return cq.queue
}

func (cq *CrawlQueue) MarkProcessed(url string) {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	delete(cq.inQueue, url)
	cq.processed[url] = true
}
