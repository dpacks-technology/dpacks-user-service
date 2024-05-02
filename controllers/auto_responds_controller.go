package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/validators"
	"fmt"
	"net/http"
	"strconv"
	_ "strconv"
	"strings"
	_ "strings"

	"github.com/gin-gonic/gin"
)

// AddAutoRespond handles POST  - CREATE
func AddAutoRespond(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		webId := c.Param("webId")

		// get the JSON data
		var autorespond models.AutoRespond
		if err := c.ShouldBindJSON(&autorespond); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		autorespond.WebID = webId

		// Validate the autorespond data
		if err := validators.ValidateMessage(autorespond, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the autorespond
		query := "INSERT INTO automated_messages (message, trigger, last_updated, status, webid) VALUES ($1, $2, $3, $4, $5)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("Error preparing statement: %s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(autorespond.Message, autorespond.Trigger, autorespond.LastUpdated, autorespond.Status, autorespond.WebID)
		if err != nil {
			fmt.Printf("Error executing statement: %s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Message added successfully"})
	}
}

func GetAutoResponds(db *sql.DB) gin.HandlerFunc {
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
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM automated_messages WHERE webid = $1 ORDER BY id LIMIT $2 OFFSET $3"
		args = append(args, c.Param("webId"), countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE id = $4 AND webid = $1 ORDER BY id LIMIT $2 OFFSET $3"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE message LIKE $4 AND webid = $1 ORDER BY CASE WHEN message = $4 THEN 1 ELSE 2 END, id LIMIT $2 OFFSET $3"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE trigger LIKE $4 AND webid = $1 ORDER BY CASE WHEN trigger = $4 THEN 1 ELSE 2 END, id LIMIT $2 OFFSET $3"
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

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into AutoRespondModel structs
		var autoresponds []models.AutoRespond

		for rows.Next() {
			var autorespond models.AutoRespond
			if err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status, &autorespond.WebID); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			autoresponds = append(autoresponds, autorespond)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, autoresponds)

	}
}

func GetAutoRespondsById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		idStr := c.Param("id")

		// get webid parameter
		webId := c.Param("webId")

		// Convert the id parameter to an integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
			return
		}

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM automated_messages WHERE id = $1 AND webid = $2", id, webId)

		// Create a AutoRespondModel to hold the data
		var autorespond models.AutoRespond

		// Scan the row data into the WebpageModel
		err = row.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status, &autorespond.WebID)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the AutoRespond as JSON
		c.JSON(http.StatusOK, autorespond)

	}
}

// GetAutoRespondByStatusCount handles  - READ
func GetAutoRespondsByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		webid := c.Param("webId")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1"
		args = append(args, webid)
		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1 AND status = 1"
			args = append(args)
		case "0":
			query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1 AND status = 0"
			args = append(args)
		default:
			query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1"
			args = append(args)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE id = $2 AND webid =  $1AND status IN ($3)"
				args = append(args, val, 1)
			case "message":
				query = "SELECT COUNT(*) FROM automated_messages WHERE message LIKE $2 AND webid = $1 AND status IN ($3)"
				args = append(args, escapedVal, 1)
			case "trigger":
				query = "SELECT COUNT(*) FROM automated_messages WHERE trigger LIKE $2 AND webid = $1 AND status IN ($3)"
				args = append(args, escapedVal, 1)
			}
		}

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing statement"})
			return
		}

		// Execute the prepared statement with bound parameters
		var count int
		err = stmt.QueryRow(args...).Scan(&count)

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Close the statement
		defer stmt.Close()

		// Return all webpages as JSON
		c.JSON(http.StatusOK, count)

	}
}

// GetAutoRespondsByStatus handles  - READ
func GetAutoRespondsByStatus(db *sql.DB) gin.HandlerFunc {
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
		webid := c.Param("webId")

		var args []interface{}
		var query string

		query = "SELECT * FROM automated_messages WHERE webid = $1 ORDER BY id LIMIT $2 OFFSET $3"
		args = append(args, webid, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM automated_messages WHERE webid = $1 AND status IN ($4) ORDER BY id LIMIT $2 OFFSET $3"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM automated_messages WHERE webid = $1 AND status IN ($4) ORDER BY id LIMIT $2 OFFSET $3"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE webid = $1 AND id = $5 AND status IN ($4) ORDER BY id LIMIT $2 OFFSET $3"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE webid = $1 AND message LIKE $5 AND status IN ($4) ORDER BY id LIMIT $2 OFFSET $3"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE webid = $1 AND trigger LIKE $5 AND status IN ($4) ORDER BY id LIMIT $2 OFFSET $3"
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

		// Iterate over the rows and scan them into AutoRespondModel structs
		var autoresponds []models.AutoRespond

		for rows.Next() {
			var autorespond models.AutoRespond
			if err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status, &autorespond.WebID); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			autoresponds = append(autoresponds, autorespond)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all AutoRespond as JSON
		c.JSON(http.StatusOK, autoresponds)

	}
}

// -----------------------------------------------------------------------------------------------------------------------//
// GetAutoRespondsByDatetime handles  - READ
func GetAutoRespondsByDatetime(db *sql.DB) gin.HandlerFunc {
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
		webid := c.Param("webId")

		var args []interface{}

		// Query the database for records based on pagination and webid
		query := "SELECT * FROM automated_messages WHERE webid = $1 ORDER BY id LIMIT $2 OFFSET $3"
		args = append(args, webid, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM automated_messages WHERE webid = $1 AND last_updated BETWEEN $4 AND $5 ORDER BY id LIMIT $2 OFFSET $3"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE webid = $1 AND id = $6 AND last_updated BETWEEN $4 AND $5 ORDER BY id LIMIT $2 OFFSET $3"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE webid = $1 AND message LIKE $6 AND last_updated BETWEEN $4 AND $5 ORDER BY CASE WHEN message = $6 THEN 1 ELSE 2 END, id LIMIT $2 OFFSET $3"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE webid = $1 AND trigger LIKE $6 AND last_updated BETWEEN $4 AND $5 ORDER BY CASE WHEN trigger = $6 THEN 1 ELSE 2 END, id LIMIT $2 OFFSET $3"
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
		var autoresponds []models.AutoRespond

		for rows.Next() {
			var autorespond models.AutoRespond
			if err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status, &autorespond.WebID); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			autoresponds = append(autoresponds, autorespond)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, autoresponds)

	}
}

// GetAutoRespondsByDatetimeCount  READ
func GetAutoRespondsByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")
		webid := c.Param("webId")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1"

		args = append(args, webid)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1 AND last_updated BETWEEN $2 AND $3"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1 AND id = $4 AND last_updated BETWEEN $2 AND $3"
				args = append(args, val)
			case "message":
				query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1 AND message LIKE $4 AND last_updated BETWEEN $2 AND $3"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT COUNT(*) FROM automated_messages WHERE webid = $1 AND trigger LIKE $4 AND last_updated BETWEEN $2 AND $3"
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

func GetAutoRespondsCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		webid := c.Param("webId")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination and webid
		query := "SELECT COUNT(*) FROM automated_messages WHERE webid = $1"
		args = append(args, webid)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE id = $2 AND webid = $1"
				args = append(args, val, webid)
			case "message":
				query = "SELECT COUNT(*) FROM automated_messages WHERE message LIKE $2 AND webid = $1"
				args = append(args, escapedVal, webid)
			case "trigger":
				query = "SELECT COUNT(*) FROM automated_messages WHERE trigger LIKE $2 AND webid = $1"
				args = append(args, escapedVal, webid)
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

// EditAutoResponds  UPDATE
func EditAutoResponds(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")
		webId := c.Param("webId")

		// get the JSON data - only the name
		var autorespond models.AutoRespond
		if err := c.ShouldBindJSON(&autorespond); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateMessage(autorespond, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage in the database
		_, err := db.Exec("UPDATE automated_messages SET message = $1 ,trigger = $2 WHERE id = $3 AND webid = $4", autorespond.Message, autorespond.Trigger, id, webId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Message/ updated successfully"})

	}
}

// DeleteAutoRespondsID DELETE
func DeleteAutoRespondsID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")
		webId := c.Param("webId")

		// query to delete the webpage
		query := "DELETE FROM automated_messages WHERE id = $1 AND web_id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(id, webId)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})

	}
}

// DeleteAutoRespondsByIDBulk  DELETE
func DeleteAutoRespondsByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")
		webId := c.Param("webId")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM automated_messages WHERE id = $1 AND webid = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(id, webId)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Message bulk deleted successfully"})

	}
}

// UpdateAutoRespondsStatusBulk h UPDATE
func UpdateAutoRespondsStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")
		webId := c.Param("webId")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var autoresponds models.AutoRespond
		if err := c.ShouldBindJSON(&autoresponds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage status in the database
		for _, id := range ids {

			query := "UPDATE automated_messages SET status = $1 WHERE id = $2 AND webid = $3"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(autoresponds.Status, id, webId)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Message status updated successfully"})

	}
}

// UpdateAutoRespondsStatus  UPDATE
func UpdateAutoRespondsStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id and webId parameters
		id := c.Param("id")
		webId := c.Param("webId")

		// get the JSON data - only the status

		var autoresponds models.AutoRespond
		if err := c.ShouldBindJSON(&autoresponds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// set the WebId field from the URL parameter
		autoresponds.WebID = webId

		// query to update the webpage status
		query := "UPDATE automated_messages SET status = $1, webid = $2 WHERE id = $3"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(autoresponds.Status, autoresponds.WebID, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Message status updated successfully"})

	}
}

func GetAutoRespondsByWebID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		webID := c.Param("webId") // Extract webID from the URL path

		// Build the query to retrieve all autoresponds for the webID
		query := "SELECT id, message, trigger, last_updated, status FROM automated_messages WHERE webid = $1"
		args := []interface{}{webID} // Bind the webID parameter

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing the database statement"})
			return
		}
		defer stmt.Close()

		// Execute the prepared statement
		rows, err := stmt.Query(args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing the database query"})
			return
		}
		defer rows.Close() // Close rows after use

		// Scan rows into AutoRespond structs
		var autoresponds []models.AutoRespond

		for rows.Next() {
			var autorespond models.AutoRespond
			err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			autoresponds = append(autoresponds, autorespond)
		}

		// Check for errors after iterating through rows
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all autoresponds as JSON
		c.JSON(http.StatusOK, autoresponds)
	}
}
