package middleware

import (
	"time"

	"backend/logger"

	"github.com/gin-gonic/gin"
)

// RequestLogger replaces Gin's default logger with structured zerolog output.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID, _ := c.Get("requestId")

		evt := logger.Log.Info()
		if status >= 500 {
			evt = logger.Log.Error()
		} else if status >= 400 {
			evt = logger.Log.Warn()
		}

		evt.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Interface("requestId", requestID).
			Int("bodySize", c.Writer.Size()).
			Msg("request")
	}
}
