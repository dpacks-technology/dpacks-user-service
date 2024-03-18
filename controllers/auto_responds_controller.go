package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAutoResponds function
func GetAutoResponds(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM automated_messages")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into AutoRespond structs
		var autoResponds []models.AutoRespond

		for rows.Next() {
			var autoRespond models.AutoRespond
			if err := rows.Scan(&autoRespond.ID, &autoRespond.Message, &autoRespond.Trigger, &autoRespond.IsActive, &autoRespond.LastUpdated); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			autoResponds = append(autoResponds, autoRespond)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all autoResponds as JSON
		c.JSON(http.StatusOK, autoResponds)

	}
}
