package main

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter enforces a sliding-window request limit per client.
type RateLimiter struct {
	maxRequests int
	timeWindow  time.Duration
	mu          sync.Mutex
	clients     map[string][]time.Time
}

// NewRateLimiter creates a rate limiter that allows at most maxRequests within timeWindow.
func NewRateLimiter(maxRequests int, timeWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
		clients:     make(map[string][]time.Time),
	}
}

// IsRequestAllowed returns true and records the request if the client has not exceeded their limit.
func (rl *RateLimiter) IsRequestAllowed(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.timeWindow)

	// Drop expired timestamps — only requests inside the sliding window count.
	// Entries are appended in order, so we scan from the front and keep the rest.
	var timestamps []time.Time
	for _, ts := range rl.clients[clientID] {
		if !ts.Before(windowStart) {
			timestamps = append(timestamps, ts)
		}
	}

	// Limit decision: deny if the client already has maxRequests in the window.
	if len(timestamps) >= rl.maxRequests {
		rl.clients[clientID] = timestamps
		return false
	}

	// Allowed — record this request and persist the pruned list.
	timestamps = append(timestamps, now)
	rl.clients[clientID] = timestamps
	return true
}

func main() {
	rl := NewRateLimiter(2, 10*time.Second)

	fmt.Println(rl.IsRequestAllowed("user_123")) // true  — 1st request
	fmt.Println(rl.IsRequestAllowed("user_123")) // true  — 2nd request
	fmt.Println(rl.IsRequestAllowed("user_123")) // false — limit reached
	fmt.Println(rl.IsRequestAllowed("user_456")) // true  — different client
}
