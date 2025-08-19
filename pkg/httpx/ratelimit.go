package httpx

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// ipRateLimiter implements a simple token-bucket per IP.
type ipRateLimiter struct {
	mu       sync.Mutex
	rpm      float64
	capacity float64
	buckets  map[string]*ipBucket
}

type ipBucket struct {
	tokens   float64
	lastFill time.Time
}

func newIPRateLimiter(requestsPerMinute int) *ipRateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = 60
	}

	return &ipRateLimiter{
		rpm:      float64(requestsPerMinute),
		capacity: float64(requestsPerMinute),
		buckets:  make(map[string]*ipBucket),
	}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[ip]
	if !ok {
		l.buckets[ip] = &ipBucket{tokens: l.capacity - 1, lastFill: time.Now()}
		return true
	}

	elapsed := time.Since(b.lastFill).Minutes()
	if elapsed > 0 {
		b.tokens += l.rpm * elapsed
		if b.tokens > l.capacity {
			b.tokens = l.capacity
		}
		b.lastFill = time.Now()
	}

	if b.tokens < 1 {
		return false
	}

	b.tokens -= 1
	return true
}

// RateLimitMiddleware returns an Echo middleware using a per-IP token bucket.
func RateLimitMiddleware(requestsPerMinute int) echo.MiddlewareFunc {
	limiter := newIPRateLimiter(requestsPerMinute)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			if !limiter.allow(ip) {
				return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "rate limit exceeded"})
			}

			return next(c)
		}
	}
}
