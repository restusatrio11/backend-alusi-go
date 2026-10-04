package middleware

import (
	"net/http"
	"sync"
	"time"

	"backend-alusi-go/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips   map[string]*visitor
	mu    sync.RWMutex
	r     rate.Limit
	b     int
	ttl   time.Duration
}

// NewIPRateLimiter creates a new rate limiter instance with background cleanup
func NewIPRateLimiter(r rate.Limit, b int, ttl time.Duration) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*visitor),
		r:   r,
		b:   b,
		ttl: ttl,
	}

	go i.cleanupVisitors()
	return i
}

func (i *IPRateLimiter) getVisitor(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.ips[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (i *IPRateLimiter) cleanupVisitors() {
	for {
		time.Sleep(1 * time.Minute)
		i.mu.Lock()
		for ip, v := range i.ips {
			if time.Since(v.lastSeen) > i.ttl {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimit middleware enforces rate limiting per client IP
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.getVisitor(ip)

		if !l.Allow() {
			response.Error(
				c,
				http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				"Terlalu banyak permintaan. Silakan tunggu beberapa saat.",
				nil,
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
