package server

import (
	"loadbalancer/internal/balancer"
	"loadbalancer/internal/proxy"
	"loadbalancer/internal/ratelimiter"
	"loadbalancer/pkg/logger"
	"net/http"
)

func StartHTTPServer(port string, rr *balancer.RoundRobin, rl *ratelimiter.Limiter) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ratelimiter.GetIP(r)
		if !rl.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			logger.Log().Printf("Rate limit exceeded for IP %s", ip)
			return
		}

		target := rr.NextBackend()
		proxy, err := proxy.NewReverseProxy(target)
		if err != nil {
			http.Error(w, "Internal proxy error", http.StatusInternalServerError)
			logger.Log().Printf("Proxy creation error: %v", err)
			return
		}

		logger.Log().Printf("Proxying request from %s to %s", ip, target)
		proxy.ServeHTTP(w, r)
	})

	logger.Log().Printf("Listening on port %s", port)
	return http.ListenAndServe(":"+port, handler)
}
