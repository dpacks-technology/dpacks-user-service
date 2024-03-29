package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/validators"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddWebPage handles POST /api/web/webpages - CREATE
func AddAutomatedMessage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var AutoRespond models.AutoRespond
		if err := c.ShouldBindJSON(&AutoRespond); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateRespond(AutoRespond, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the webpage
		query := "INSERT INTO automated_messages (message, trigger, is_active, last_updated) VALUES ($1, $2, $3, $4)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		result, err := stmt.Exec(AutoRespond.Message, AutoRespond.Trigger, AutoRespond.IsActive, AutoRespond.LastUpdated)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Check if the row was inserted successfully
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Set the ID of the newly createdAutoRespond
		AutoRespond.ID = int(lastInsertId)

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"AutoRespond": AutoRespond})
	}
}

// GetWebPages handles GET /api/web/pages/ - READ
func GetAutomatedMessage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Param("page")
		count := c.Param("count")

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		offset := (pageInt - 1) * countInt

		key := c.Query("key")
		val := c.Query("val")

		// Validate the AutoRespond model before querying the database
		var AutoRespond models.AutoRespond
		if err := c.ShouldBindQuery(&AutoRespond); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validators.ValidateRespond(AutoRespond, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		query := "SELECT * FROM automated_messages WHERE is_active = true ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE is_active = true AND id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val, escapedVal)
			case "message":
				query = "SELECT * FROM automated_messages WHERE is_active = true AND message LIKE $3 ORDER BY CASE WHEN message = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE is_active = true AND trigger LIKE $3 ORDER BY CASE WHEN trigger = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			}
		}

		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		defer stmt.Close()

		rows, err := stmt.Query(args...)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		defer rows.Close()

		var autoResponses []models.AutoRespond

		for rows.Next() {
			var autoRespond models.AutoRespond
			if err := rows.Scan(&autoRespond.ID, &autoRespond.Message, &autoRespond.Trigger, &autoRespond.IsActive, &autoRespond.LastUpdated); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			autoResponses = append(autoResponses, autoRespond)
		}

		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		c.JSON(http.StatusOK, autoResponses)
	}
}

/*
// GetWebPageById handles GET /api/web/webpages/:id - READ
func GetAutomatedMessageById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM automated_messages WHERE id = $1", id)

		// Create an AutoRespond to hold the data
		var autoRespond models.AutoRespond

		// Scan the row data into the AutoRespond
		err := row.Scan(&autoRespond.ID, &autoRespond.Message, &autoRespond.Trigger, &autoRespond.IsActive, &autoRespond.LastUpdated)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the AutoRespond as JSON
		c.JSON(http.StatusOK, autoRespond)

	}
}

// GetWebPagesByStatusCount handles GET /api/web/webpages/status/:status/count - READ
func GetAutomatedMessageStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM automated_messages"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM automated_messages WHERE is_active IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM automated_messages WHERE is_active IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE id = $2 AND is_active IN ($1)"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM automated_messages WHERE message LIKE $2 AND is_active IN ($1)"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT COUNT(*) FROM automated_messages WHERE trigger LIKE $2 AND is_active IN ($1)"
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

// GetWebPagesByStatus handles GET /api/web/webpages/status/:status - READ
func GetAutomatedMessageByStatus(db *sql.DB) gin.HandlerFunc {
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

		query = "SELECT * FROM automated_messages ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM automated_messages WHERE is_active IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM automated_messages WHERE is_active IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE is_active IN ($3) AND id = $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val, escapedVal)
			case "name":
				query = "SELECT * FROM automated_messages WHERE is_active IN ($3) AND message LIKE $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val, escapedVal)
			case "path":
				query = "SELECT * FROM automated_messages WHERE is_active IN ($3) AND trigger LIKE $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val, escapedVal)
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

		// Iterate over the rows and scan them into AutoRespond structs
		var autoResponses []models.AutoRespond

		for rows.Next() {
			var autoRespond models.AutoRespond
			if err := rows.Scan(&autoRespond.ID, &autoRespond.Message, &autoRespond.Trigger, &autoRespond.IsActive, &autoRespond.LastUpdated); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows fromthe database"})
				return
			}
			autoResponses = append(autoResponses, autoRespond)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all autoResponses as JSON
		c.JSON(http.StatusOK, autoResponses)

	}
}

// GetWebPagesByDatetime handles GET /api/web/webpages/datetime/:count/:page - READ
func GetAutomatedMessageByDatetime(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM automated_messages ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM automated_messages WHERE date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE id = $5 AND date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE message LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN message = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE trigger LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN trigger = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into AutoRespond structs
		var automatedMessages []models.AutoRespond

		for rows.Next() {
			var automatedMessage models.AutoRespond
			if err := rows.Scan(&automatedMessage.ID, &automatedMessage.Message, &automatedMessage.Trigger, &automatedMessage.IsActive, &automatedMessage.LastUpdated); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			automatedMessages = append(automatedMessages, automatedMessage)
		}

		//this runs only when loop didn'twork
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all automated messages as JSON
		c.JSON(http.StatusOK, automatedMessages)
	}
}

// GetWebPagesByDatetimeCount handles GET /api/web/webpages/datetime/count - READ
func GetAutomatedMessageByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM webpages"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM webpages WHERE date_created BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM webpages WHERE id = $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM webpages WHERE name LIKE $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT COUNT(*) FROM webpages WHERE path LIKE $3 AND date_created BETWEEN $1 AND $2"
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

func GetAutomatedMessageCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM webpages"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM webpages WHERE id = $1"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM webpages WHERE name LIKE $1"
				args = append(args, escapedVal)
			case "path":
				query = "SELECT COUNT(*) FROM webpages WHERE path LIKE $1"
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

// EditWebPage handles PUT /api/web/webpages/:id - UPDATE
func EditAutomatedMessage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var webpage models.WebpageModel
		if err := c.ShouldBindJSON(&webpage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateName(webpage, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage in the database
		_, err := db.Exec("UPDATE webpages SET name = $1 WHERE id = $2", webpage.Name, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})

	}
}

// DeleteWebPageByID handles DELETE /api/web/webpages/:id - DELETE
func DeleteAutomatedMessageByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM webpages WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Webpage deleted successfully"})

	}
}

// DeleteWebPageByIDBulk handles DELETE /api/web/webpages/bulk/:id - DELETE
func DeleteAutomatedMessageByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM webpages WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Webpage bulk deleted successfully"})

	}
}

// UpdateWebPageStatusBulk handles PUT /api/web/webpages/status/bulk/:id - UPDATE
func UpdateAutomatedMessageStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var webpage models.WebpageModel
		if err := c.ShouldBindJSON(&webpage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage status in the database
		for _, id := range ids {

			query := "UPDATE webpages SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(webpage.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}
}

// UpdateWebPageStatus handles PUT /api/web/webpages/status/:id - UPDATE
func UpdateAutomatedMessageStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var webpage models.WebpageModel
		if err := c.ShouldBindJSON(&webpage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE webpages SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(webpage.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}
}
*/
