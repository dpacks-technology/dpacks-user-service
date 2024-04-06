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
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query with the website ID as a parameter
		rows, err := db.Query(`
            SELECT
                EXTRACT(DOW FROM sessionstart) AS day_of_week,
                COUNT(*) AS session_count
            FROM
                public.sessions
            WHERE
                web_id = $1
            GROUP BY
                day_of_week
            ORDER BY
                day_of_week;
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
			var dayOfWeek int
			var sessionCount int
			if err := rows.Scan(&dayOfWeek, &sessionCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"day_of_week": dayOfWeek, "session_count": sessionCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}
func GetDevices(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query with the website ID as a parameter
		rows, err := db.Query(`
            SELECT
                d.devicename,
                COUNT(*) AS device_count
            FROM
                public.sessions AS s
            JOIN
                public.devices AS d ON s.deviceid = d.deviceid
            WHERE
                s.web_id IN ($1)
            GROUP BY
                s.web_id,
                d.devicename
            ORDER BY
                s.web_id,
                d.devicename;
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
			var deviceName string
			var deviceCount int
			if err := rows.Scan(&deviceName, &deviceCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"device_name": deviceName, "device_count": deviceCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}

func GetCountry(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {
		// Get the website ID from the URL params
		websiteID := c.Param("id")

		// Execute the SQL query with the website ID as a parameter
		rows, err := db.Query(`
            SELECT
                c.countrycode,
                COUNT(DISTINCT s.ipaddress) AS user_count
            FROM
                public.sessions AS s
            JOIN
                public.countries AS c ON s.countrycode = c.countrycode
            WHERE
                s.web_id = $1
            GROUP BY
                s.web_id,
                c.countrycode
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
			var countryCode string
			var userCount int
			if err := rows.Scan(&countryCode, &userCount); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
				return
			}
			// Append the result to the results slice
			results = append(results, gin.H{"country_code": countryCode, "user_count": userCount})
		}

		// Return the results as JSON
		c.JSON(http.StatusOK, results)
	}
}
