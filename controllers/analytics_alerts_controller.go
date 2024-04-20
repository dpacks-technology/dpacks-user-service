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
