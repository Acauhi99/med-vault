package httpx

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const maxBuckets = 10000

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    int
	window  time.Duration
	stop    chan struct{}
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		window:  window,
		stop:    make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		now := time.Now()

		rl.mu.Lock()
		b, ok := rl.buckets[ip]
		if !ok || now.Sub(b.lastReset) >= rl.window {
			if !ok && len(rl.buckets) >= maxBuckets {
				rl.mu.Unlock()
				http.Error(w, `{"error":{"code":"RATE_LIMIT_EXCEEDED","message":"Too many requests"}}`, http.StatusTooManyRequests)
				return
			}
			rl.buckets[ip] = &bucket{tokens: rl.rate - 1, lastReset: now}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if b.tokens <= 0 {
			elapsed := now.Sub(b.lastReset)
			retryAfter := int((rl.window - elapsed).Seconds()) + 1
			rl.mu.Unlock()
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			http.Error(w, `{"error":{"code":"RATE_LIMIT_EXCEEDED","message":"Too many requests"}}`, http.StatusTooManyRequests)
			return
		}

		b.tokens--
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) Close() {
	close(rl.stop)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-rl.stop:
			return
		case now := <-ticker.C:
			rl.mu.Lock()
			for ip, b := range rl.buckets {
				if now.Sub(b.lastReset) > 10*time.Minute {
					delete(rl.buckets, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

func extractIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
