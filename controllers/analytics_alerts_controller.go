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
		query := "SELECT * FROM webpages ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM webpages WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM webpages WHERE name LIKE $3 ORDER BY CASE WHEN name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT * FROM webpages WHERE path LIKE $3 ORDER BY CASE WHEN path = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
		var webpages []models.WebpageModel

		for rows.Next() {
			var webpage models.WebpageModel
			if err := rows.Scan(&webpage.ID, &webpage.Name, &webpage.WebID, &webpage.Path, &webpage.Status, &webpage.DateCreated); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			webpages = append(webpages, webpage)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, webpages)

	}

}
