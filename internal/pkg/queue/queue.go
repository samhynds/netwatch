package queue

import (
	"log"
	"net/http"
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

type Queue[T string | ProcessedQueueItem] struct {
	Queue   chan T
	inQueue map[string]bool
	mu      sync.RWMutex
	qType   QueueType
	parent  *Queues
}

type ProcessedQueueItem struct {
	URL       string
	Timestamp time.Time
	Content   map[string]string
	Links     []string
	Body      string
	Headers   http.Header
}

type Queues struct {
	Crawl     Queue[string]
	Cooldown  Queue[string]
	Recrawl   Queue[string]
	Transport Queue[ProcessedQueueItem]
	processed map[string]bool
	mu        sync.RWMutex
}

func NewQueue(config *config.Config) *Queues {
	q := &Queues{
		processed: make(map[string]bool),
	}

	q.Crawl = Queue[string]{
		Queue:   make(chan string, config.Config.Queue.Capacity),
		inQueue: make(map[string]bool),
		qType:   CrawlQueue,
		parent:  q,
	}

	q.Cooldown = Queue[string]{
		Queue:   make(chan string, config.Config.Queue.Capacity),
		inQueue: make(map[string]bool),
		qType:   CooldownQueue,
		parent:  q,
	}

	q.Recrawl = Queue[string]{
		Queue:   make(chan string, len(config.Sites)),
		inQueue: make(map[string]bool),
		qType:   RecrawlQueue,
		parent:  q,
	}

	q.Transport = Queue[ProcessedQueueItem]{
		Queue:   make(chan ProcessedQueueItem, config.Config.Queue.Capacity),
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
			q.Recrawl.Queue <- site.URL
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
	case ProcessedQueueItem:
		url = v.URL
	}

	if q.inQueue[url] {
		log.Println("Already in queue:", url)
		return false
	}

	// Only check processed URLs for crawl queue
	if q.qType == CrawlQueue {
		q.parent.mu.RLock()
		processed := q.parent.processed[url]
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
	case ProcessedQueueItem:
		url = v.URL
	}

	if q.qType == CrawlQueue {
		q.parent.mu.Lock()
		defer q.parent.mu.Unlock()

		delete(q.inQueue, url)
		q.parent.processed[url] = true
	} else {
		log.Println("Cannot mark non-crawl queue items as processed", url)
	}
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
