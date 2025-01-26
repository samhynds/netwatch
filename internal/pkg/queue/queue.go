package queue

import (
	"log"
	"net/http"
	"net/url"
	"netwatch/internal/pkg/config"
	"sync"
	"time"
)

type QueueType int

const (
	CrawlQueue QueueType = iota
	CooldownQueue
	RecrawlQueue
	TransportQueue
)

type Queue[T string | CrawledQueueItem | ScheduledQueueItem] struct {
	Queue   chan T
	inQueue map[string]bool
	mu      sync.RWMutex
	qType   QueueType
	parent  *Queues
}

type CrawledQueueItem struct {
	URL       string
	Timestamp time.Time
	Content   map[string]string
	Links     []string
	Body      string
	Headers   http.Header
}

type ScheduledQueueItem struct {
	URL           string
	TimeCrawlable time.Time
}

type ProcessedHost struct {
	count int // Number of requests made to this host within the last window
}

type Queues struct {
	Crawl          Queue[string]
	Cooldown       Queue[ScheduledQueueItem]
	Recrawl        Queue[ScheduledQueueItem]
	Transport      Queue[CrawledQueueItem]
	processedUrls  map[string]bool
	processedHosts map[string]ProcessedHost
	mu             sync.RWMutex
}

func NewQueue(config *config.Config) *Queues {
	q := &Queues{
		processedUrls:  make(map[string]bool),
		processedHosts: make(map[string]ProcessedHost),
	}

	q.Crawl = Queue[string]{
		Queue:   make(chan string, config.Config.Queue.Capacity),
		inQueue: make(map[string]bool),
		qType:   CrawlQueue,
		parent:  q,
	}

	q.Cooldown = Queue[ScheduledQueueItem]{
		Queue:   make(chan ScheduledQueueItem, config.Config.Queue.Capacity),
		inQueue: make(map[string]bool),
		qType:   CooldownQueue,
		parent:  q,
	}

	q.Recrawl = Queue[ScheduledQueueItem]{
		Queue:   make(chan ScheduledQueueItem, len(config.Sites)),
		inQueue: make(map[string]bool),
		qType:   RecrawlQueue,
		parent:  q,
	}

	q.Transport = Queue[CrawledQueueItem]{
		Queue:   make(chan CrawledQueueItem, config.Config.Queue.Capacity),
		inQueue: make(map[string]bool),
		qType:   TransportQueue,
		parent:  q,
	}

	return q
}

func (q *Queues) InitPopulation(config *config.Config) {
	for _, site := range config.Sites {
		q.Crawl.Queue <- site.URL
	}

	if config.Config.Recrawl.Enabled {
		for _, site := range config.Sites {
			q.Recrawl.Queue <- ScheduledQueueItem{
				URL:           site.URL,
				TimeCrawlable: time.Now().Add(time.Duration(config.Config.Recrawl.Interval) * time.Second),
			}
		}
	}
}

func (q *Queue[T]) Add(item T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	var url string

	switch v := any(item).(type) {
	case string:
		url = v
	case CrawledQueueItem:
	case ScheduledQueueItem:
		url = v.URL
	}

	if q.inQueue[url] {
		log.Println("Already in queue:", url)
		return false
	}

	// Only check processed URLs for crawl queue
	if q.qType == CrawlQueue {
		q.parent.mu.RLock()
		processed := q.parent.processedUrls[url]
		q.parent.mu.RUnlock()

		if processed {
			log.Println("Already processed:", url)
			return false
		}
	}

	q.inQueue[url] = true
	q.Queue <- item
	return true
}

func (q *Queue[T]) AddMultiple(items []T) []T {
	added := make([]T, 0, len(items))
	for _, item := range items {
		if q.Add(item) {
			added = append(added, item)
		}
	}
	return added
}

func (q *Queue[T]) Get() <-chan T {
	return q.Queue
}

func (q *Queue[T]) MarkProcessed(item T) {
	var url string

	switch v := any(item).(type) {
	case string:
		url = v
	case CrawledQueueItem:
		url = v.URL
	}

	q.parent.mu.Lock()
	defer q.parent.mu.Unlock()

	delete(q.inQueue, url)
	q.parent.processedUrls[url] = true

	host, err := GetHostFromURL(url)
	if err != nil {
		log.Println("Error getting host for", url)
		return
	}

	if existingHost, exists := q.parent.processedHosts[host]; exists {
		q.parent.processedHosts[host] = ProcessedHost{
			count: existingHost.count + 1,
		}
	} else {
		q.parent.processedHosts[host] = ProcessedHost{
			count: 1,
		}
	}
}

func GetHostFromURL(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}

func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	close(q.Queue)
}

func (q *Queues) Close() {
	q.Crawl.Close()
	q.Cooldown.Close()
	q.Recrawl.Close()
}
