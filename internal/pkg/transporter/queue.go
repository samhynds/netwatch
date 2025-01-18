package transporter

import (
	"net/http"
	"sync"
	"time"
)

type TransportQueue struct {
	queue chan TransportQueueItem
	mu    sync.RWMutex
}

type TransportQueueItem struct {
	URL       string
	Timestamp time.Time
	Content   map[string]string
	Links     []string
	Body      string
	Headers   http.Header
}

func NewTransportQueue(size int) *TransportQueue {
	return &TransportQueue{
		queue: make(chan TransportQueueItem, size),
	}
}

func (tq *TransportQueue) Add(item TransportQueueItem) bool {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.queue <- item
	return true
}

func (tq *TransportQueue) Get() <-chan TransportQueueItem {
	return tq.queue
}

func (tq *TransportQueue) Close() {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	close(tq.queue)
}
