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
		query := "SELECT * FROM visitor_info ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {

			case "device":
				query = "SELECT * FROM visitor_info WHERE device LIKE $3 ORDER BY CASE WHEN device = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "country":
				query = "SELECT * FROM visitor_info WHERE country LIKE $3 ORDER BY CASE WHEN country = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
			if err := rows.Scan(&visitorInfo.ID, &visitorInfo.IpAddres, &visitorInfo.Device, &visitorInfo.Country, &visitorInfo.Source); err != nil {
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

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM visitor_info"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM visitor_info WHERE id = $1"
				args = append(args, val)
			case "device":
				query = "SELECT COUNT(*) FROM visitor_info WHERE device LIKE $1"
				args = append(args, escapedVal)
			case "country":
				query = "SELECT COUNT(*) FROM visitor_info WHERE country LIKE $1"
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

//func GetVisitorInfoByDatetime(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get page id parameter
//		page := c.Param("page")
//
//		// get count parameter
//		count := c.Param("count")
//
//		// Convert page and count to integers
//		pageInt, err := strconv.Atoi(page)
//		if err != nil {
//			// Handle error
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		countInt, err := strconv.Atoi(count)
//		if err != nil {
//			// Handle error
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Calculate offset
//		offset := (pageInt - 1) * countInt
//
//		// get query parameters
//		start := c.Query("start")
//		end := c.Query("end")
//		key := c.Query("key")
//		val := c.Query("val")
//
//		var args []interface{}
//
//		// Query the database for records based on pagination
//		query := "SELECT * FROM visitor_info ORDER BY id LIMIT $1 OFFSET $2"
//		args = append(args, countInt, offset)
//
//		if start != "" && end != "" && val != "null" && key != "null" {
//			query = "SELECT * FROM webpages WHERE date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
//			args = append(args, start, end)
//		}
//
//		if val != "" && key != "" {
//			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
//			switch key {
//			case "id":
//				query = "SELECT * FROM webpages WHERE id = $5 AND date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
//				args = append(args, val)
//			case "name":
//				query = "SELECT * FROM webpages WHERE name LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
//				args = append(args, escapedVal)
//			case "path":
//				query = "SELECT * FROM webpages WHERE path LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN path = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
//				args = append(args, escapedVal)
//			}
//		}
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Execute the prepared statement with bound parameters
//		rows, err := stmt.Query(args...)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		//close the rows when the surrounding function returns(handler function)
//		defer rows.Close()
//
//		// Iterate over the rows and scan them into WebpageModel structs
//		var webpages []models.WebpageModel
//
//		for rows.Next() {
//			var webpage models.WebpageModel
//			if err := rows.Scan(&webpage.ID, &webpage.Name, &webpage.WebID, &webpage.Path, &webpage.Status, &webpage.DateCreated); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
//				return
//			}
//			webpages = append(webpages, webpage)
//		}
//
//		//this runs only when loop didn't work
//		if err := rows.Err(); err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
//			return
//		}
//
//		// Return all webpages as JSON
//		c.JSON(http.StatusOK, webpages)
//
//	}
//
//}
