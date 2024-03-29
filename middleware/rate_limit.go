package middleware

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	db       *sql.DB
}

func NewRateLimit(db *sql.DB) (*RateLimit, error) {
	rl := &RateLimit{limiters: make(map[string]*rate.Limiter), db: db}
	// Fetch initial limits from database on creation
	err := rl.updateLimitsFromDatabase()
	if err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RateLimit) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var path string
		path = c.FullPath() // Use FullPath method to get the complete path

		rl.mu.Lock()
		limiter, found := rl.limiters[path]
		rl.mu.Unlock()

		if !found {
			// Handle missing limit (fetch from database or use default)
			if rl.db != nil {
				err := rl.updateLimitsFromDatabase()
				if err != nil {
					// Log the error and proceed without rate limiting
					c.Next()
					return
				}
				rl.mu.Lock()
				limiter, found = rl.limiters[path]
				rl.mu.Unlock()
			}
		}

		if !found {
			// No limit found in database (handle accordingly)
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

func (rl *RateLimit) updateLimitsFromDatabase() error {
	rows, err := rl.db.Query("SELECT path, ratelimit FROM endpoint_ratelimits")
	if err != nil {
		return err
	}
	defer rows.Close()

	limits := make(map[string]*rate.Limiter)
	for rows.Next() {
		var path string
		var limit int
		err := rows.Scan(&path, &limit)
		if err != nil {
			return err
		}
		limits[path] = rate.NewLimiter(rate.Every(time.Minute), limit) // Adjust rate as needed
	}

	// Update the rate limiters map
	rl.mu.Lock()
	rl.limiters = limits
	rl.mu.Unlock()

	return nil
}
