package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAnalyticalAlerts function
func GetAnalyticalAlerts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM useralerts")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into UserAlerts structs
		var userAlerts []models.UserAlerts

		for rows.Next() {
			var userAlert models.UserAlerts
			if err := rows.Scan(&userAlert.AlertID, &userAlert.UserID, &userAlert.UserEmail, &userAlert.AlertThreshold, &userAlert.AlertSubject, &userAlert.AlertContent, &userAlert.WhenAlertRequired, &userAlert.ReminderOption, &userAlert.CustomReminderDate); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			userAlerts = append(userAlerts, userAlert)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all userAlerts as JSON
		c.JSON(http.StatusOK, userAlerts)

	}
}
