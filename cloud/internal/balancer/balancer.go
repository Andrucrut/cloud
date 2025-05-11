package balancer

type Balancer interface {
	NextBackend() string
}

type RoundRobin struct {
	backends []string
	current  int
}

func NewRoundRobin(backends []string) *RoundRobin {
	return &RoundRobin{backends: backends}
}

func (rr *RoundRobin) NextBackend() string {
	backend := rr.backends[rr.current]
	rr.current = (rr.current + 1) % len(rr.backends)
	return backend
}
