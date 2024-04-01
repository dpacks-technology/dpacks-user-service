package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/validators"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM useralerts WHERE website_id = $3 ORDER BY alertid LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset, id)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM useralerts WHERE alertid = $3 AND website_id = $4 ORDER BY alertid LIMIT $1 OFFSET $2"
				args = append(args, val, id)
			case "alertthreshold":
				query = "SELECT * FROM useralerts WHERE alertthreshold LIKE $3 AND website_id = $4 ORDER BY CASE WHEN alertthreshold = $3 THEN 1 ELSE 2 END, alertid LIMIT $1 OFFSET $2"
				args = append(args, escapedVal, id)
			case "whenalertrequired":
				query = "SELECT * FROM useralerts WHERE whenalertrequired LIKE $3 AND website_id = $4 ORDER BY CASE WHEN whenalertrequired = $3 THEN 1 ELSE 2 END, alertid LIMIT $1 OFFSET $2"
				args = append(args, escapedVal, id)
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

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into WebpageModel structs
		var alerts []models.UserAlertsModel

		for rows.Next() {
			var alert models.UserAlertsModel
			if err := rows.Scan(&alert.AlertID, &alert.WebsiteeId, &alert.AlertSubject, alert.AlertThreshold, &alert.UserID, &alert.UserEmail, &alert.Status, &alert.AlertContent, &alert.CustomReminderDate, &alert.RepeatOn, &alert.WhenAlertRequired); err != nil {
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
		query := "SELECT COUNT(*) FROM useralerts WHERE website_id = $2"
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
