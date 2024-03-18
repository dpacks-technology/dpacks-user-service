package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWebPages handles GET /api/web/pages/ - READ
func GetWebPages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM webpages")

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
			if err := rows.Scan(&webpage.ID, &webpage.Name, &webpage.WebID, &webpage.Path); err != nil {
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
