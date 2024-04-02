package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func CreateNewAlert(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// get the JSON data
		var alert models.UserAlertsModel
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the webpage
		query := "INSERT INTO useralert (id, user_id, alert_threshold, alert_subject,alert_content,when_alert_required,repeat_on,custom_reminder_date,status,website_id) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,1,$9)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(alert.AlertID, alert.UserID, alert.AlertThreshold, alert.AlertSubject, alert.AlertContent, alert.WhenAlertRequired, alert.RepeatOn, alert.CustomReminderDate, alert.Status, alert.WebsiteeId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Alert set Succesfully"})

	}
}

func GetAllAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		//get id parameter
		id := c.Param("id")

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		//get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM useralerts WHERE website_id=$3 ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset, id)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM useralerts WHERE alert_id = $3  ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
				fmt.Printf("%s\n %d\n %d\n", args, val)

			case "alertthreshold":
				query = "SELECT * FROM useralerts WHERE alert_threshold LIKE $3  ORDER BY CASE WHEN alert_threshold = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "whenalertrequired":
				query = "SELECT * FROM useralerts WHERE when_alert_required LIKE $3  ORDER BY CASE WHEN when_alert_required = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {

			fmt.Printf("%s\n", err)
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer rows.Close()

		//close the rows when the surrounding function returns(handler function)

		// Iterate over the rows and scan them into WebpageModel structs
		var alerts []models.UserAlertsModel
		for rows.Next() {
			var alert models.UserAlertsModel
			if err := rows.Scan(&alert.AlertID, &alert.UserID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.RepeatOn, &alert.CustomReminderDate, &alert.Status, &alert.WebsiteeId); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			fmt.Printf("%s\n", alert)
			alerts = append(alerts, alert)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, alerts)

	}

}

func GetAlertsCount(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var count int
		id := c.Param("id")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM useralerts WHERE website_id = $1"
		args = append(args, id)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM useralerts WHERE alertid = $1 AND website_id = $2"
				args = append(args, val, id)
			case "alertthreshold":
				query = "SELECT COUNT(*) FROM useralerts WHERE alertthreshold LIKE $1 AND website_id = $2"
				args = append(args, escapedVal, id)
			case "whenalertrequired":
				query = "SELECT COUNT(*) FROM useralerts WHERE whenalertrequired LIKE $1 AND website_id = $2"
				args = append(args, escapedVal, id)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		err = stmt.QueryRow(args...).Scan(&count)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Close the statement
		defer stmt.Close()

		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)

	}

}

func GetAlertbyId(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM useralerts WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var alert models.UserAlertsModel

		// Scan the row data into the WebpageModel
		err := row.Scan(&alert.AlertID, &alert.UserID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.RepeatOn, &alert.CustomReminderDate, &alert.Status, &alert.WebsiteeId)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, alert)

	}

}

func GetAlertsByStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get page id parameter
		page := c.Param("page")

		// get count parameter
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			// Handle error
			fmt.Printf("%s\n", err)
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT * FROM useralerts ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM useralerts WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM useralerts WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM useralerts WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM useralerts WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "alertthreshold":
				query = "SELECT * FROM useralerts WHERE alertthreshold LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "whenalertrequired":
				query = "SELECT * FROM useralerts WHERE whenalertrequired LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var alerts []models.UserAlertsModel

		for rows.Next() {
			var alert models.UserAlertsModel
			if err := rows.Scan(&alert.AlertID, &alert.UserID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.RepeatOn, &alert.CustomReminderDate, &alert.Status, &alert.WebsiteeId); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			alerts = append(alerts, alert)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, alerts)

	}

}
