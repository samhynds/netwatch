package transporterqueue

import (
	"net/http"
	"sync"
	"time"
)

type Queue struct {
	queue chan QueueItem
	mu    sync.RWMutex
}

type QueueItem struct {
	URL       string
	Timestamp time.Time
	Content   map[string]string
	Links     []string
	Body      string
	Headers   http.Header
}

func NewQueue(size int) *Queue {
	return &Queue{
		queue: make(chan QueueItem, size),
	}
}

func (tq *Queue) Add(item QueueItem) bool {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.queue <- item
	return true
}

func (tq *Queue) Get() <-chan QueueItem {
	return tq.queue
}

func (tq *Queue) Close() {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	close(tq.queue)
}
