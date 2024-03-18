package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSubscriptionPlans func
func GetSubscriptionPlans(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM subscriptionplans")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into SubscriptionPlans structs
		var subscriptionPlans []models.SubscriptionPlans

		for rows.Next() {
			var subscriptionPlan models.SubscriptionPlans
			if err := rows.Scan(&subscriptionPlan.PlanID, &subscriptionPlan.PlanName, &subscriptionPlan.Desc, &subscriptionPlan.Price, &subscriptionPlan.Duration, &subscriptionPlan.Features); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			subscriptionPlans = append(subscriptionPlans, subscriptionPlan)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all subscriptionPlans as JSON
		c.JSON(http.StatusOK, subscriptionPlans)

	}
}
