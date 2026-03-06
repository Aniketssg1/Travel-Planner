package middleware

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

const requestIDHeader = "X-Request-ID"

// RequestID generates a unique ULID for every request, stores it in
// the Gin context, and sets the X-Request-ID response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.New(rand.NewSource(int64(ulid.Now()))), 0)).String()
		c.Set("requestId", id)
		c.Header(requestIDHeader, id)
		c.Next()
	}
}
