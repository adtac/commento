package main

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	tokens    float64
	maxTokens float64
	refillRate float64 // tokens per second
	lastTime  time.Time
}

var limiter *rateLimiter

func rateLimiterInit() {
	limiter = &rateLimiter{
		visitors: make(map[string]*visitor),
	}

	// Periodically clean up stale entries
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limiter.cleanup()
		}
	}()
}

func (rl *rateLimiter) getVisitor(key string, maxTokens float64, refillRate float64) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists {
		v = &visitor{
			tokens:    maxTokens,
			maxTokens: maxTokens,
			refillRate: refillRate,
			lastTime:  time.Now(),
		}
		rl.visitors[key] = v
		return v
	}

	return v
}

func (v *visitor) allow() bool {
	now := time.Now()
	elapsed := now.Sub(v.lastTime).Seconds()
	v.lastTime = now

	v.tokens += elapsed * v.refillRate
	if v.tokens > v.maxTokens {
		v.tokens = v.maxTokens
	}

	if v.tokens < 1 {
		return false
	}

	v.tokens--
	return true
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	for key, v := range rl.visitors {
		if v.lastTime.Before(cutoff) {
			delete(rl.visitors, key)
		}
	}
}

// rateLimitMiddleware applies per-IP rate limiting.
// maxPerMinute is the maximum requests allowed per minute.
func rateLimitMiddleware(maxPerMinute float64, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter == nil {
			next(w, r)
			return
		}

		ip := getIp(r)
		v := limiter.getVisitor("ip:"+ip, maxPerMinute, maxPerMinute/60.0)

		if !v.allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			bodyMarshal(w, response{"success": false, "message": "Rate limit exceeded. Please try again later."})
			return
		}

		next(w, r)
	}
}
