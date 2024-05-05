package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//GetStorageByID handles GET /api/web/storage/:id - READ

func GetStorageByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT SUM(size) FROM data_packets WHERE site = $1", id)

		// Create a WebpageModel to hold the data
		var Storage models.DataPackets

		// Scan the row data into the WebpageModel
		err := row.Scan(&Storage.Size)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, Storage)

	}
}
