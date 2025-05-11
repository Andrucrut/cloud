package ratelimiter

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type Limiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func NewLimiter(refillRate int, capacity int) *Limiter {
	return &Limiter{
		limiters: make(map[string]*rate.Limiter),
		r:        rate.Limit(refillRate),
		b:        capacity,
	}
}

func (l *Limiter) getLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	limiter, exists := l.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.limiters[ip] = limiter
	}
	return limiter
}

func (l *Limiter) Allow(ip string) bool {
	return l.getLimiter(ip).Allow()
}

func GetIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
