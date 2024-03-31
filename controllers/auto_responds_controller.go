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

// AddAutoRespond handles POST /api/web/webpages - CREATE
func AddAutoRespond(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var autorespond models.AutoRespond
		if err := c.ShouldBindJSON(&autorespond); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateMessage(autorespond, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the webpage
		query := "INSERT INTO automated_messages (message, trigger, last_updated, status) VALUES ($1, $2, $3, $4)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("Error preparing statement: %s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		result, err := stmt.Exec(autorespond.Message, autorespond.Trigger, autorespond.LastUpdated, autorespond.Status)
		if err != nil {
			fmt.Printf("Error executing statement: %s\n", err)
			return
		}

		// Check if the insert was successful
		lastID, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("Error getting last insert ID: %s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Message added successfully", "id": lastID})
	}
}

// GetWebPages handles GET /api/web/pages/ - READ
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
		query := "SELECT * FROM automated_messages ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE message LIKE $3 ORDER BY CASE WHEN message = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE trigger LIKE $3 ORDER BY CASE WHEN trigger = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into WebpageModel structs
		var autoresponds []models.AutoRespond

		for rows.Next() {
			var autorespond models.AutoRespond
			if err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status); err != nil {
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

// GetWebPageById handles GET /api/web/webpages/:id - READ
func GetAutoRespondsById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM automated_messages WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var autorespond models.AutoRespond

		// Scan the row data into the WebpageModel
		err := row.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, autorespond)

	}
}

// GetWebPagesByStatusCount handles GET /api/web/webpages/status/:status/count - READ
func GetAutoRespondsByStatusCount(db *sql.DB) gin.HandlerFunc {
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
			query = "SELECT COUNT(*) FROM automated_messages WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM automated_messages WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE id = $2 AND status IN ($1)"
				args = append(args, val)
			case "message":
				query = "SELECT COUNT(*) FROM automated_messages WHERE message LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT COUNT(*) FROM automated_messages WHERE trigger LIKE $2 AND status IN ($1)"
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

		var args []interface{}
		var query string

		query = "SELECT * FROM automated_messages ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM automated_messages WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM automated_messages WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM automated_messages WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE message LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE trigger LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
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
			if err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status); err != nil {
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

// -----------------------------------------------------------------------------------------------------------------------//
// GetWebPagesByDatetime handles GET /api/web/webpages/datetime/:count/:page - READ
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

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT * FROM automated_messages ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM automated_messages WHERE last_updated BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM automated_messages WHERE id = $5 AND last_updated BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "message":
				query = "SELECT * FROM automated_messages WHERE name LIKE $5 AND last_updated BETWEEN $3 AND $4 ORDER BY CASE WHEN message = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT * FROM automated_messages WHERE path LIKE $5 AND last_updated BETWEEN $3 AND $4 ORDER BY CASE WHEN trigger = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
			if err := rows.Scan(&autorespond.ID, &autorespond.Message, &autorespond.Trigger, &autorespond.LastUpdated, &autorespond.Status); err != nil {
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

// GetWebPagesByDatetimeCount handles GET /api/web/webpages/datetime/count - READ
func GetAutoRespondsByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM automated_messages"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM automated_messages WHERE last_updated BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE id = $3 AND last_updated BETWEEN $1 AND $2"
				args = append(args, val)
			case "message":
				query = "SELECT COUNT(*) FROM automated_messages WHERE message LIKE $3 AND last_updated BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT COUNT(*) FROM automated_messages WHERE trigger LIKE $3 AND last_updated BETWEEN $1 AND $2"
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
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM automated_messages"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM automated_messages WHERE id = $1"
				args = append(args, val)
			case "message":
				query = "SELECT COUNT(*) FROM automated_messages WHERE message LIKE $1"
				args = append(args, escapedVal)
			case "trigger":
				query = "SELECT COUNT(*) FROM automated_messages WHERE trigger LIKE $1"
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
func EditAutoResponds(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

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
		_, err := db.Exec("UPDATE automated_messages SET message = $1 WHERE id = $2", autorespond.Message, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})

	}
}

// DeleteWebPageByID handles DELETE /api/web/webpages/:id - DELETE
func DeleteAutoRespondsID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM automated_messages WHERE id = $1"

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
func DeleteAutoRespondsByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM automated_messages WHERE id = $1"

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
func UpdateAutoRespondsStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

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

			query := "UPDATE automated_messages SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(autoresponds.Status, id)
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
func UpdateAutoRespondsStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status

		var autoresponds models.AutoRespond
		if err := c.ShouldBindJSON(&autoresponds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE automated_messages SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(autoresponds.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}
}
