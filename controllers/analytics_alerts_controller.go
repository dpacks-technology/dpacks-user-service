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
func GetSource(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query
		rows, err := db.Query(`
            SELECT
                src.type AS user_source,
                COUNT(*) AS user_count
            FROM
                public.sessions AS s
            JOIN
                public.source AS src ON s.source_id = src.id
            WHERE
                s.web_id = $1
            GROUP BY
                s.web_id,
                src.type
            ORDER BY
                s.web_id,
                user_count DESC;
        `, websiteID)
		if err != nil {
			// Handle any errors
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()

		// Iterate over the result set
		var results []gin.H
		for rows.Next() {
			var userSource string
			var userCount int
			if err := rows.Scan(&userSource, &userCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"user_source": userSource, "user_count": userCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}

func GetSessions(db *sql.DB) gin.HandlerFunc {

}

//func GetDevices(db *sql.DB) gin.HandlerFunc {
//
//}
//
//func GetCountry(db *sql.DB) gin.HandlerFunc {
//
//}
