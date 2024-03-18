package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetTemplates function
func GetTemplates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM templates")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into Template structs
		var templates []models.Template

		for rows.Next() {
			var template models.Template
			if err := rows.Scan(&template.TemplId, &template.TemplName, &template.TemplDescription, &template.TemplCategory, &template.MainTemplFile, &template.ThmbnlTemplFile, &template.TemplDevpName, &template.UserID, &template.TemplDevpDescription); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all templates as JSON
		c.JSON(http.StatusOK, templates)

	}
}
