package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetWebPages handles GET /api/web/pages/ - READ
func GetWebPages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// Query the database for records based on pagination
		rows, err := db.Query("SELECT * FROM webpages LIMIT $1 OFFSET $2", countInt, offset)

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
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

func GetWebPagesCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM webpages").Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)
	}
}
