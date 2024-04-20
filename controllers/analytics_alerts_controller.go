package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetVisitorInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM visitor_info WHERE webid = $3 ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset, id) // Added 'id' to the args slice for SQL query

		if val != "" && key != "" {
			switch key {
			case "id":
				// Already filtering by ID, so no need to change the query
			case "device":
				query = "SELECT * FROM visitor_info WHERE webid = $3 AND device LIKE $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal) // Adjust the args slice accordingly
			case "country":
				query = "SELECT * FROM visitor_info WHERE webid = $3 AND country LIKE $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var visitorsInfo []models.VisitorInfo

		for rows.Next() {
			var visitorInfo models.VisitorInfo
			if err := rows.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source, &visitorInfo.VisitedDate, &visitorInfo.WebId); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			visitorsInfo = append(visitorsInfo, visitorInfo)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, visitorsInfo)
	}
}

func GetVisitorInfoCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the visitor ID from the URL parameter
		id := c.Param("id")

		var count int

		// Get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Start building the query
		query := "SELECT COUNT(*) FROM visitor_info WHERE webid = $1"
		args = append(args, id)

		if val != "" && key != "" {
			// Depending on the key, modify the query to include an additional filter
			switch key {
			case "id":
				// No additional filter needed if key is 'id' since we're already filtering by id
			case "device":
				query += " AND device LIKE $2"
				args = append(args, escapedVal)
			case "country":
				query += " AND country LIKE $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the SQL statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing query"})
			return
		}
		defer stmt.Close()

		// Execute the query
		err = stmt.QueryRow(args...).Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing query"})
			return
		}

		// Return the count as JSON
		c.JSON(http.StatusOK, count)
	}
}

func GetVisitorInfoByDatetime(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM visitor_info ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM visitor_info WHERE visited_time BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM visitor_info WHERE id = $5 AND visited_time BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM visitor_info WHERE device LIKE $5 AND visited_time BETWEEN $3 AND $4 ORDER BY CASE WHEN device = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT * FROM visitor_info WHERE country LIKE $5 AND visited_time BETWEEN $3 AND $4 ORDER BY CASE WHEN country = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var visitorsInfo []models.VisitorInfo

		for rows.Next() {
			var visitorInfo models.VisitorInfo
			if err := rows.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source, &visitorInfo.VisitedDate); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			visitorsInfo = append(visitorsInfo, visitorInfo)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, visitorsInfo)

	}

}

func GetVisitorByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM visitor_info"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM visitor_info WHERE visited_time BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM visitor_info WHERE id = $3 AND visited_time BETWEEN $1 AND $2"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM visitor_info WHERE device LIKE $3 AND visited_time BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT COUNT(*) FROM visitor_info WHERE country LIKE $3 AND visited_time BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		var count int
		err = stmt.QueryRow(args...).Scan(&count)

		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Close the statement
		defer stmt.Close()

		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)

	}

}

func GetVisitorInfoById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM visitor_info WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var visitorInfo models.VisitorInfo

		// Scan the row data into the WebpageModel
		err := row.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source, &visitorInfo.VisitedDate)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, visitorInfo)

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
