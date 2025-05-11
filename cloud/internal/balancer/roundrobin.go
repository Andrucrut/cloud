package balancer

import (
	"sync"
)

type RoundRobin struct {
	backends []string
	mu       sync.Mutex
	current  int
}

func NewRoundRobin(backends []string) *RoundRobin {
	return &RoundRobin{backends: backends}
}

func (r *RoundRobin) GetNextBackend() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	backend := r.backends[r.current]
	r.current = (r.current + 1) % len(r.backends)
	return backend
}
