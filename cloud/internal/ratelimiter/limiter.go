package ratelimiter

import (
	"sync"
	"time"
)

type Limiter struct {
	limits     map[string]*clientLimiter
	mu         sync.Mutex
	refillRate int
	capacity   int
}

type clientLimiter struct {
	tokens     int
	lastRefill time.Time
}

func NewLimiter(refillRate, capacity int) *Limiter {
	return &Limiter{
		limits:     make(map[string]*clientLimiter),
		refillRate: refillRate,
		capacity:   capacity,
	}
}

func (l *Limiter) AllowRequest(clientID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.limits[clientID]; !exists {
		l.limits[clientID] = &clientLimiter{tokens: l.capacity, lastRefill: time.Now()}
	}

	client := l.limits[clientID]
	now := time.Now()

	elapsed := now.Sub(client.lastRefill)
	refillTokens := int(elapsed.Seconds()) * l.refillRate

	client.tokens = min(client.tokens+refillTokens, l.capacity)
	client.lastRefill = now

	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
