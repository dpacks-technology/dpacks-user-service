package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/validators"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateNewAlert(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// get the JSON data
		var webpage models.WebpageModel
		if err := c.ShouldBindJSON(&webpage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateName(webpage, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the webpage
		query := "INSERT INTO webpages (name, webid, path, status) VALUES ($1, $2, $3, $4)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(webpage.Name, webpage.WebID, webpage.Path, 1)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Webpage added successfully"})

	}
}
