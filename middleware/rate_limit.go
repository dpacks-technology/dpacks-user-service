package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimit struct {
	limiters map[string]*rate.Limiter // Map endpoint paths to rate limiters
}

func NewRateLimit(limits map[string]int) *rateLimit {
	rl := &rateLimit{limiters: make(map[string]*rate.Limiter)}
	for path, limit := range limits {
		rl.limiters[path] = rate.NewLimiter(rate.Every(time.Minute), limit) // Adjust rate as needed
	}
	return rl
}

func (rl *rateLimit) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		limiter, ok := rl.limiters[path]
		if !ok {
			// Handle cases where a path doesn't have a specific limit set
			// You might choose to use a default limiter or allow the request
			c.Next()
			return
		}

		// Check allowance without waiting (using a non-blocking approach)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded for this resource. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
