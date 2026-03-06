package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout wraps every request context with a deadline.
// Default: 30 seconds. Configurable via REQUEST_TIMEOUT_SECONDS.
func Timeout() gin.HandlerFunc {
	timeout := 30 * time.Second
	if v := os.Getenv("REQUEST_TIMEOUT_SECONDS"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			timeout = time.Duration(secs) * time.Second
		}
	}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timed out",
				"code":  "TIMEOUT",
			})
			c.Abort()
		}
	}
}
