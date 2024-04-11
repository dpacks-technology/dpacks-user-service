package middleware

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware is a middleware function to authenticate requests using API key and client ID
func AuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("API-Key")     // Assuming API key is passed in the "X-API-Key" header
		clientID := c.GetHeader("Client-ID") // Assuming Client ID is passed in the "X-Client-ID" header

		// Check if both API key and client ID are provided
		if apiKey == "" || clientID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key or Client ID missing"})
			return
		}

		// Query the database to check if the provided API key and client ID are valid
		var userID int
		err := db.QueryRow("SELECT user_id FROM api_subscribers WHERE key = $1 AND client_id = $2", apiKey, clientID).Scan(&userID)
		if err != nil {
			// Invalid API key or client ID
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key or Client ID"})
			return
		}

		// Set the user ID in the context for further processing
		c.Set("userID", userID)
		// Respond with a success message with the set user ID
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated", "userID": userID})

		// Continue to the next handler
		c.Next()
	}
}
