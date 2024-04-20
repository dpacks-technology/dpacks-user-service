package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSubscriotionID handles GET /api/web/subscription/:id - READ
func GetSubscriptionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM subscription WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var Subscription models.SubscriptionModel

		// Scan the row data into the WebpageModel
		err := row.Scan(&Subscription.ID, &Subscription.ProjectID, &Subscription.PlanID, &Subscription.PlanName, &Subscription.Amount)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, Subscription)

	}
}

// DeleteSubscriptionByID handles DELETE /api/web/subscription/:id - DELETE
func DeleteSubscriptionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Execute the SQL DELETE statement
		_, err := db.Exec("DELETE FROM subscription WHERE id = $1", id)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting row from the database"})
			return
		}

		// Return a 204 status
		c.JSON(http.StatusNoContent, nil)

	}
}
