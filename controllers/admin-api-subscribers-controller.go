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

// AddSubscribers
func AddSubscribers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var subscriber models.ApiSubscriber
		if err := c.ShouldBindJSON(&subscriber); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the subscriber data
		if err := validators.ValidateUserId(subscriber, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		subscriber.ClientID = generateUniqueID()
		subscriber.Key = generateUniqueID()

		// query to insert the keypair
		query := "INSERT INTO api_subscribers (user_id, client_id, key) VALUES ($1, $2, $3)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(subscriber.UserID, subscriber.ClientID, subscriber.Key)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Subscriber added successfully"})

	}
}

// GetApiSubscribers
func GetApiSubscribers(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM api_subscribers ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM api_subscribers WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "user_id":
				query = "SELECT * FROM api_subscribers WHERE user_id LIKE $3 ORDER BY CASE WHEN user_id = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "client_id":
				query = "SELECT * FROM api_subscribers WHERE client_id LIKE $3 ORDER BY CASE WHEN client_id = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into KeypairModel structs
		var subscribers []models.ApiSubscriber

		for rows.Next() {
			var subscriber models.ApiSubscriber
			if err := rows.Scan(&subscriber.ID, &subscriber.UserID, &subscriber.ClientID, &subscriber.Key, &subscriber.CreatedOn); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			subscribers = append(subscribers, subscriber)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all subscriber as JSON
		c.JSON(http.StatusOK, subscribers)

	}
}

// GetApiSubscriberById
func GetApiSubscriberById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM api_subscribers WHERE id = $1", id)

		// Create a KeypairModel to hold the data
		var subscriber models.ApiSubscriber

		// Scan the row data into the KeypairModel
		err := row.Scan(&subscriber.ID, &subscriber.UserID, &subscriber.ClientID, &subscriber.Key, &subscriber.CreatedOn)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the Keypair as JSON
		c.JSON(http.StatusOK, subscriber)

	}
}

// GetApiSubscribersByDatetime
func GetApiSubscribersByDatetime(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM api_subscribers ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM api_subscribers WHERE created_on BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM api_subscribers WHERE id = $5 AND created_on BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "user_id":
				query = "SELECT * FROM api_subscribers WHERE user_id LIKE $5 AND created_on BETWEEN $3 AND $4 ORDER BY CASE WHEN user_id = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "client_id":
				query = "SELECT * FROM api_subscribers WHERE client_id LIKE $5 AND created_on BETWEEN $3 AND $4 ORDER BY CASE WHEN client_id = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into KeypairModel structs
		var subscribers []models.ApiSubscriber

		for rows.Next() {
			var subscriber models.ApiSubscriber
			if err := rows.Scan(&subscriber.ID, &subscriber.UserID, &subscriber.ClientID, &subscriber.Key, &subscriber.CreatedOn); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			subscribers = append(subscribers, subscriber)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all keypair as JSON
		c.JSON(http.StatusOK, subscribers)

	}
}

// GetApiSubscribersByDatetimeCount
func GetApiSubscribersByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM api_subscribers"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM api_subscribers WHERE created_on BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM api_subscribers WHERE id = $3 AND created_on BETWEEN $1 AND $2"
				args = append(args, val)
			case "user_id":
				query = "SELECT COUNT(*) FROM api_subscribers WHERE user_id LIKE $3 AND created_on BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "client_id":
				query = "SELECT COUNT(*) FROM api_subscribers WHERE client_id LIKE $3 AND created_on BETWEEN $1 AND $2"
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

		c.JSON(http.StatusOK, count)

	}
}

// GetApiSubscribersCount
func GetApiSubscribersCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM api_subscribers"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM api_subscribers WHERE id = $1"
				args = append(args, val)
			case "client_id":
				query = "SELECT COUNT(*) FROM api_subscribers WHERE client_id LIKE $1"
				args = append(args, escapedVal)
			case "user_id":
				query = "SELECT COUNT(*) FROM api_subscribers WHERE user_id LIKE $1"
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

		c.JSON(http.StatusOK, count)

	}
}

// RegenerateKey
func RegenerateKey(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var subscriber models.ApiSubscriber

		subscriber.ClientID = generateUniqueID()
		subscriber.Key = generateUniqueID()

		// Update the keypair in the database
		_, err := db.Exec("UPDATE api_subscribers SET client_id = $1,key = $2  WHERE id = $3", subscriber.ClientID, subscriber.Key, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Keypair updated successfully"})

	}
}

// DeleteApiSubscriberByID
func DeleteApiSubscriberByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the Keypair
		query := "DELETE FROM api_subscribers WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Keypair deleted successfully"})

	}
}

// DeleteApiSubscriberByIDBulk
func DeleteApiSubscriberByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the keypair from the database
		for _, id := range ids {
			// query to delete the Keypair
			query := "DELETE FROM api_subscribers WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Keypair bulk deleted successfully"})

	}
}
