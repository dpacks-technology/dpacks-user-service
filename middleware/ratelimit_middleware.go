package middleware

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// creating ratelimit struct
type RateLimit struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	db       *sql.DB
}

// creating  ratelimit instance
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

// creating limit function
func (rl *RateLimit) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var path string
		path = c.FullPath() // Use FullPath method to get the complete path
		// Log the path
		log.Println("path: ", path)

		rl.mu.Lock()
		limiter, found := rl.limiters[path]
		rl.mu.Unlock()

		// If the path is not found in the map, return a 404 Not Found error
		if !found {
			//log the message
			log.Println("No such Endpoint")
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No such Endpoint, Please try again later.",
			})
			return
		}

		// Check the status of the endpoint
		var limitModel models.Endpoint
		limitModel.Status, _ = rl.statusCheck(path)

		// Check if the endpoint is enabled or disabled
		if limitModel.Status == 0 {
			//log the message
			log.Println("Endpoint is disabled")

			// Return a message to display this endpoint is not working
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Endpoint Currently Disabled. Please try again later.",
			})

			return
		} else {
			//log the message
			log.Println("Endpoint is enabled")

			// Check allowance without waiting (using a non-blocking approach)
			if !limiter.Allow() {

				//log the message
				log.Println("Rate limit exceeded")

				// If the rate limit is exceeded, return a 429 Too Many Requests error
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded for this resource. Please try again later.",
				})
				return
			}
			// Call the next handler if the rate limit is not exceeded
			c.Next()
		}

	}
}

// function to Update the rate limiters from the database
func (rl *RateLimit) updateLimitsFromDatabase() error {
	// Query the database for the rate limits
	rows, err := rl.db.Query("SELECT path, ratelimit FROM api_endpoints")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a map to store the rate limiters
	limits := make(map[string]*rate.Limiter)

	// Iterate over the rows and scan them into the Endpoint struct
	for rows.Next() {
		//using model
		var limitModel models.Endpoint

		err := rows.Scan(&limitModel.Path, &limitModel.Limit)
		if err != nil {
			return err
		}
		// Create a new rate limiter for the path
		limits[limitModel.Path] = rate.NewLimiter(rate.Every(time.Minute), limitModel.Limit) // Adjust rate as needed
	}

	// Update the rate limiters map
	rl.mu.Lock()
	rl.limiters = limits
	rl.mu.Unlock()

	return nil
}

// function to check the status of the endpoint
func (rl *RateLimit) statusCheck(path string) (int, error) {
	// Check if the endpoint is enabled or disabled
	// Query the database for the status of the endpoint
	query := "SELECT status FROM api_endpoints WHERE path = $1"

	// Prepare the statement
	stmt, err := rl.db.Prepare(query)
	if err != nil {
		fmt.Printf("%s\n", "Error preparing the query")
		return 0, err
	}

	// Close the statement when the surrounding function returns (handler function)
	defer stmt.Close()

	// Execute the statement
	row, err := stmt.Query(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 0, err
	}

	// Close the rows when the surrounding function returns (handler function)
	defer row.Close()

	// Create a model to store the status of the endpoint
	var limitModel models.Endpoint

	// Iterate over the rows and scan them into the Endpoint struct
	for row.Next() {
		err := row.Scan(&limitModel.Status)
		if err != nil {
			fmt.Printf("%s\n", err)
			return 0, err
		}
	}

	// Return the status of the endpoint
	return limitModel.Status, err

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
