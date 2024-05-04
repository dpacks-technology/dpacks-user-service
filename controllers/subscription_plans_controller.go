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

		// check count first and return empty array if no records
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM subscriptionplans").Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting count of subscription plans from the database"})
			return
		}

		if count == 0 {
			c.JSON(http.StatusOK, []models.SubscriptionPlans{})
			return
		}

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM subscriptionplans ORDER BY plan_id")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting subscription plans from the database"})
			return
		}

		// Create a slice of SubscriptionPlans
		var subscriptionPlans []models.SubscriptionPlans

		// Iterate over the rows, adding each SubscriptionPlan to the slice
		for rows.Next() {
			var subscriptionPlan models.SubscriptionPlans
			if err := rows.Scan(&subscriptionPlan.PlanID, &subscriptionPlan.PlanName, &subscriptionPlan.Desc, &subscriptionPlan.Duration, &subscriptionPlan.Features, &subscriptionPlan.MonthlyPrice, &subscriptionPlan.AnnualPrice, &subscriptionPlan.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
				return
			}
			subscriptionPlans = append(subscriptionPlans, subscriptionPlan)
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, subscriptionPlans)

	}
}
