package ratelimiter

import (
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	mutex    sync.Mutex
	maxTotal int
	window   time.Duration
	requests []time.Time
}

func NewRateLimiter(maxTotal int, window int) *RateLimiter {
	return &RateLimiter{
		maxTotal: maxTotal,
		window:   time.Duration(window) * time.Second,
		requests: make([]time.Time, 0, maxTotal),
	}
}

func (rl *RateLimiter) Allow() (bool, time.Time) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove outdated requests
	for len(rl.requests) > 0 && rl.requests[0].Before(windowStart) {
		rl.requests = rl.requests[1:]
	}

	// Check if we've reached the limit
	if len(rl.requests) >= rl.maxTotal {
		nextAllowedTime := rl.requests[0].Add(rl.window)
		return false, nextAllowedTime
	}

	// Add the new request
	rl.requests = append(rl.requests, now)
	log.Printf("Request Count: %d", len(rl.requests))
	log.Println("Requests: ", rl.requests)
	return true, now
}
