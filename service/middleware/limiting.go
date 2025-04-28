package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter manages rate limiting based on client IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter initializes a new RateLimiter instance
func NewRateLimiter(rps int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps), // Requests per second
		burst:    burst,           // Burst size
	}
}

// getLimiter retrieves or creates a rate limiter for a specific client IP
func (r *RateLimiter) getLimiter(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiter, exists := r.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(r.rate, r.burst)
		r.limiters[ip] = limiter

		// Clean up unused limiters after a set duration (e.g., 10 minutes)
		go func(ip string) {
			time.Sleep(10 * time.Minute)
			r.mu.Lock()
			delete(r.limiters, ip)
			r.mu.Unlock()
		}(ip)
	}

	return limiter
}

// RateLimitMiddleware applies rate limiting to incoming requests
func RateLimitMiddleware(rps int, burst int) gin.HandlerFunc {
	rl := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		limiter := rl.getLimiter(clientIP)

		if !limiter.Allow() {
			// Rate limit exceeded, return HTTP 429
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
