package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

// GetAllWebContents function
func GetAllWebContents(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Return a message to display this endpoint is working
		c.JSON(200, gin.H{"message": "GetAllWebContents endpoint is working!!!!!!!!!!"})

	}
}
