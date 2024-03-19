package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// // GetAnalyticalAlerts function
//
//	func GetAnalyticalAlerts(db *sql.DB) gin.HandlerFunc {
//		return func(c *gin.Context) {
//
//			// Query the database for all records
//			rows, err := db.Query("SELECT * FROM useralerts")
//
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
//				return
//			}
//
//			//close the rows when the surrounding function returns(handler function)
//			defer rows.Close()
//
//			// Iterate over the rows and scan them into UserAlerts structs
//			var userAlerts []models.UserAlerts
//
//			for rows.Next() {
//				var userAlert models.UserAlerts
//				if err := rows.Scan(&userAlert.AlertID, &userAlert.UserID, &userAlert.UserEmail, &userAlert.AlertThreshold, &userAlert.AlertSubject, &userAlert.AlertContent, &userAlert.WhenAlertRequired, &userAlert.ReminderOption, &userAlert.CustomReminderDate); err != nil {
//					fmt.Printf("%s\n", err)
//					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
//					return
//				}
//				userAlerts = append(userAlerts, userAlert)
//			}
//
//			//this runs only when loop didn't work
//			if err := rows.Err(); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
//				return
//			}
//
//			// Return all userAlerts as JSON
//			c.JSON(http.StatusOK, userAlerts)
//
//		}
//	}
//
// delete alert

// update alert

// create  new alert
func CreateAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var createAlert models.CreateNewAlert

		if err := c.BindJSON(&createAlert); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error binding JSON"})
			return
		}

		_, err := db.Exec("INSERT INTO useralerts (userid,useremail,alertthreshold,alertsubject,alertcontent,whenalertrequired,reminderoption) VALUES ($1, $2, $3, $4, $5, $6, $7)", createAlert.UserID, createAlert.UserEmail, createAlert.AlertThreshold, createAlert.AlertSubject, createAlert.AlertContent, createAlert.WhenAlertRequired, createAlert.ReminderOption)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting into the database"})
			return
		}
		// Return the example as JSON
		c.JSON(http.StatusOK, createAlert)

	}
}

func DeleteAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		//chech if this id availble in the database
		var alertID int
		err := db.QueryRow("SELECT alertid FROM useralerts WHERE alertid=$1", id).Scan(&alertID)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user id not available in the database"})
			return
		}

		result, err := db.Exec("DELETE FROM useralerts WHERE alertid=$1", id)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting from the database"})
			return
		}
		_, err = result.RowsAffected()
		if err != nil {
			fmt.Printf("%s\n", err)
			fmt.Printf("No rows affected\n")
		}
		// Return the example as JSON
		c.JSON(http.StatusOK, gin.H{"Delete": "Alert Deleted"})
	}
}

func UpdateAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var updateAlert models.CreateNewAlert

		if err := c.BindJSON(&updateAlert); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error binding JSON"})
			return
		}

		_, err := db.Exec("UPDATE useralerts SET userid=$1,useremail=$2,alertthreshold=$3, alertsubject=$4, alertcontent=$5, whenalertrequired=$6, reminderoption=$7 WHERE alertid=$8", updateAlert.UserID, updateAlert.UserEmail, updateAlert.AlertThreshold, updateAlert.AlertSubject, updateAlert.AlertContent, updateAlert.WhenAlertRequired, updateAlert.ReminderOption, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the database"})
			return
		}
		// Return the example as JSON
		c.JSON(http.StatusOK, updateAlert)

	}
}

// alert list show
func GetAlertList(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		rows, err := db.Query("SELECT alertthreshold,alertsubject,alertcontent,whenalertrequired,reminderoption,customreminderdate FROM useralerts WHERE userid = $1", id)

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()

		var alertList []models.AlertList

		for rows.Next() {
			var alert models.AlertList
			if err := rows.Scan(&alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.ReminderOption, &alert.CustomReminderDate); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			alertList = append(alertList, alert)

		}
		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all examples as JSON
		c.JSON(http.StatusOK, alertList)

	}

}

func GetDashboardData(db *sql.DB) gin.HandlerFunc {
	
}