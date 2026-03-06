package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	tokens    float64
	lastSeen  time.Time
	maxTokens float64
	rate      float64 // tokens per second
}

var (
	visitors   = make(map[string]*visitor)
	visitorsMu sync.Mutex
)

func init() {
	// Background cleanup: remove stale entries every 3 minutes.
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			visitorsMu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(visitors, ip)
				}
			}
			visitorsMu.Unlock()
		}
	}()
}

func getVisitor(ip string, rps float64) *visitor {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		v = &visitor{tokens: rps, maxTokens: rps * 2, rate: rps, lastSeen: time.Now()}
		visitors[ip] = v
		return v
	}
	return v
}

func (v *visitor) allow() bool {
	now := time.Now()
	elapsed := now.Sub(v.lastSeen).Seconds()
	v.lastSeen = now
	v.tokens += elapsed * v.rate
	if v.tokens > v.maxTokens {
		v.tokens = v.maxTokens
	}
	if v.tokens < 1 {
		return false
	}
	v.tokens--
	return true
}

// RateLimiter returns a per-IP token-bucket rate limiter middleware.
// Default: 100 req/s per IP, configurable via RATE_LIMIT_RPS env var.
func RateLimiter() gin.HandlerFunc {
	rps := 100.0
	if v := os.Getenv("RATE_LIMIT_RPS"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil && parsed > 0 {
			rps = parsed
		}
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		v := getVisitor(ip, rps)

		visitorsMu.Lock()
		allowed := v.allow()
		visitorsMu.Unlock()

		if !allowed {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please slow down",
				"code":  "TOO_MANY_REQUESTS",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
