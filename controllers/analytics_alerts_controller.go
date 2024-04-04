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
		var alert models.CreateNewUserAlert
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//print the model alert
		//fmt.Printf("test %s", alert)

		// query to insert the webpage
		query := "INSERT INTO useralerts (alert_threshold, alert_subject,alert_content,when_alert_required,repeat_on,website_id) VALUES ($1, $2, $3, $4,$5,$6)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fmt.Printf("test3")

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(alert.AlertThreshold, alert.AlertSubject, alert.AlertContent, alert.WhenAlertRequired, alert.RepeatOn, alert.WebsiteeId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fmt.Printf("test4")

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
		//escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT id,alert_threshold,alert_subject,alert_content,when_alert_required,repeat_on,status,website_id FROM useralerts WHERE website_id=$3 ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset, id)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM useralerts WHERE id = $3  ORDER BY id LIMIT $1 OFFSET $2"
				fmt.Printf("%s", args)
				args = append(args, val)

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
		var alerts []models.UserAlertsShow
		for rows.Next() {
			var alert models.UserAlertsShow
			if err := rows.Scan(&alert.AlertID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.RepeatOn, &alert.Status, &alert.WebsiteeId); err != nil {
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
				query = "SELECT COUNT(*) FROM useralerts WHERE id = $1 AND website_id = $2"
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
		row := db.QueryRow("SELECT id,alert_threshold,alert_subject,alert_content,when_alert_required,repeat_on,status,website_id FROM useralerts WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var alert models.UserAlertsShow

		// Scan the row data into the WebpageModel
		err := row.Scan(&alert.AlertID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.RepeatOn, &alert.Status, &alert.WebsiteeId)
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
		var alerts []models.UserAlertsShow

		for rows.Next() {
			var alert models.UserAlertsShow
			if err := rows.Scan(&alert.AlertID, &alert.AlertThreshold, &alert.AlertSubject, &alert.AlertContent, &alert.WhenAlertRequired, &alert.RepeatOn, &alert.Status, &alert.WebsiteeId); err != nil {
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
func GetAlertsByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM useralerts"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM useralerts WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM useralerts WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM useralerts WHERE id = $2 AND status IN ($1)"
				args = append(args, val)
			case "alertthreshold":
				query = "SELECT COUNT(*) FROM useralerts WHERE alertthreshold LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "whenalertrequired":
				query = "SELECT COUNT(*) FROM useralerts WHERE whenalertrequired LIKE $2 AND status IN ($1)"
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
		var count int
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

//func EditAlert(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// get the JSON data - only the name
//		var alert models.UserAlertsModel
//		if err := c.ShouldBindJSON(&alert); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		//// Validate the webpage data
//		//if err := validators.ValidateName(alert, false); err != nil {
//		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		//	return
//		//}
//
//		// Update the webpage in the database
//		_, err := db.Exec("UPDATE useralerts SET name = $1 WHERE id = $2", alert.Name, id)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})
//
//	}
//
//}

func DeleteAlertByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM useralerts WHERE id = $1"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})

	}

}

func DeleteAlertByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM useralerts WHERE id = $1"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Alert bulk deleted successfully"})

	}

}

func UpdateAlertStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var alert models.UserAlertStatus
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE useralerts SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(alert.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}

}

func UpdateAlertStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var alert models.UserAlertStatus
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage status in the database
		for _, id := range ids {

			query := "UPDATE useralerts SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(alert.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}
}

func EditAlert(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var updateAlert models.UpdateUserAlert
		if err := c.ShouldBindJSON(&updateAlert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("test %s", updateAlert)

		// Validate the webpage data
		//if err := validators.ValidateName(webpage, false); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}

		// Update the webpage in the database
		_, err := db.Exec("UPDATE useralerts SET alert_threshold=$1,alert_subject=$2,alert_content=$3,when_alert_required=$4,repeat_on=$5 where id=$6", updateAlert.AlertThreshold, updateAlert.AlertSubject, updateAlert.AlertContent, updateAlert.WhenAlertRequired, updateAlert.RepeatOn, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})

	}
}
