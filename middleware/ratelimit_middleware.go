package middleware

import (
	"database/sql"
	"dpacks-go-services-template/models"
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

	// Start a goroutine to periodically refresh the rate limits from the database
	go rl.periodicRefresh()

	return rl, nil
}

func (rl *RateLimit) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var path string
		path = c.FullPath() // Use FullPath method to get the complete path

		rl.mu.Lock()
		limiter, found := rl.limiters[path]
		rl.mu.Unlock()

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
	rows, err := rl.db.Query("SELECT path, ratelimit FROM api_endpoints")
	if err != nil {
		return err
	}
	defer rows.Close()

	limits := make(map[string]*rate.Limiter)
	for rows.Next() {

		//using model
		var limitModel models.Endpoint

		err := rows.Scan(&limitModel.Path, &limitModel.Limit)
		if err != nil {
			return err
		}
		limits[limitModel.Path] = rate.NewLimiter(rate.Every(time.Minute), limitModel.Limit) // Adjust rate as needed
	}

	// Update the rate limiters map
	rl.mu.Lock()
	rl.limiters = limits
	rl.mu.Unlock()

	return nil
}

// Periodically refresh rate limits from the database
// becouse when we update the rate limit in the database we need to update the rate limit in the instance we have in the memory
// so we need to refresh the rate limit from the database
func (rl *RateLimit) periodicRefresh() {
	// Function to periodically refresh the rate limits from the database
	// This function runs in a separate goroutine

	// Define the refresh interval (e.g., every 30 seconds
	refreshInterval := 1 * time.Minute

	// Create a ticker to trigger refresh at regular intervals
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Refresh rate limits from the database
			if err := rl.updateLimitsFromDatabase(); err != nil {

			}
		}
	}
}
