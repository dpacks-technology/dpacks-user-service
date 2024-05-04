package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSubscriptionByID handles GET /api/web/subscription/:id - READ
func GetSubscriptionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// check count first and return empty array if no records
		var count int
		err := db.QueryRow("SELECT COUNT(plan.plan_id) FROM subscription sub, subscriptionplans plan WHERE sub.plan_id = plan.plan_id AND sub.project_id = $1", id).Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting count of subscription from the database"})
			return
		}

		if count == 0 {
			c.JSON(http.StatusOK, []models.SubscriptionPlans{})
			return
		}

		// Query the database for a single record
		row := db.QueryRow("SELECT plan.plan_id, plan.plan_name, plan.description, plan.duration, plan.features, plan.monthly_price, plan.annual_price, plan.status FROM subscription sub, subscriptionplans plan WHERE sub.plan_id = plan.plan_id AND sub.project_id = $1", id)

		// Scan the row into a Subscription struct
		var SubscriptionPlan models.SubscriptionPlans
		if err := row.Scan(&SubscriptionPlan.PlanID, &SubscriptionPlan.PlanName, &SubscriptionPlan.Desc, &SubscriptionPlan.Duration, &SubscriptionPlan.Features, &SubscriptionPlan.MonthlyPrice, &SubscriptionPlan.AnnualPrice, &SubscriptionPlan.Status); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, SubscriptionPlan)

	}
}

// DeleteSubscriptionByID handles DELETE /api/web/subscription/:id - DELETE
func DeleteSubscriptionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Execute the SQL DELETE statement and call the sendEmail function

		_, err := db.Exec("DELETE FROM subscription WHERE project_id = $1", id)
		if err != nil {

			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting row from the database"})
			return
		}

		// Send an email to the user
		err = utils.SendEmail("erandi14908@gmail.com", "Subscription Cancelled", "Your subscription has been cancelled", "small")
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email"})
			return
		}

		// Return a 204 status
		c.JSON(http.StatusNoContent, nil)

	}

}
