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
		rows, err := db.Query("SELECT * FROM subscriptionplans ORDER BY plan_id")

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
			if err := rows.Scan(&subscriptionPlan.PlanID, &subscriptionPlan.PlanName, &subscriptionPlan.Desc, &subscriptionPlan.Duration, &subscriptionPlan.Features, &subscriptionPlan.MonthlyPrice, &subscriptionPlan.AnnualPrice, &subscriptionPlan.Status); err != nil {
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

// AddExample handles POST /api/example - CREATE

func AddSubscriptionplan(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Create an empty ExampleModel struct
		var transactions models.transaction

		// Bind the JSON to the ExampleModel struct
		if err := c.BindJSON(&transactions); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON"})
			return
		}

		// Insert the record into the database
		_, err := db.Exec("INSERT INTO transactions (column2, column3, new_column) VALUES ($1, $2, $3)", example.Column1, example.Column2, example.Column3)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting into the database"})
			return
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, example)

	}
}
