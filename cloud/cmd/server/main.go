package main

import (
	"loadbalancer/config"
	"loadbalancer/internal/balancer"
	"loadbalancer/internal/ratelimiter"
	"loadbalancer/internal/server"
	"loadbalancer/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig("/Users/andrey/GolandProjects/cloud/config/config.yaml")
	if err != nil {
		logger.Log().Fatal(err)
	}

	rr := balancer.NewRoundRobin(cfg.Backends)
	rl := ratelimiter.NewLimiter(cfg.RateLimit.RefillRate, cfg.RateLimit.Capacity)

	if err := server.StartHTTPServer(cfg.ListenPort, rr, rl); err != nil {
		logger.Log().Fatal(err)
	}
}
