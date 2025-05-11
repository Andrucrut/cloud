package server

import (
	"encoding/json"
	"loadbalancer/internal/balancer"
	"loadbalancer/internal/ratelimiter"
	"loadbalancer/pkg/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func StartHTTPServer(port string, rr *balancer.RoundRobin, rl *ratelimiter.Limiter) error {
	log := logger.Log()
	log.Printf("Listening on port %s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr

		if !rl.AllowRequest(clientID) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    429,
				"message": "Rate limit exceeded",
			})
			return
		}

		backend := rr.GetNextBackend()
		log.Printf("Proxying request from %s to %s\n", clientID, backend)

		target, err := url.Parse(backend)
		if err != nil {
			http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}
