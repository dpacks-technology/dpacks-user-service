package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

// GetKeyPairs function
func GetKeyPairs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		query := "SELECT * FROM api_subscribers"

		//prepare statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing the query"})
			return
		}
		//close the statement when the surrounding function returns(handler function)
		defer stmt.Close()

		//execute the statement
		rows, err := stmt.Query()
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}
		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into KeyPairs structs
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

		// Return all keypairs as JSON
		c.JSON(http.StatusOK, subscribers)

	}
}

func GetKeyPairsID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		// Create an empty ExampleModel struct
		var subscriber models.ApiSubscriber

		//query the database for the record with the given ID
		query := "SELECT * FROM api_subscribers WHERE user_id = $1"

		//prepare statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing the query"})
			return
		}

		//close the statement when the surrounding function returns(handler function)
		defer stmt.Close()

		//execute the statement
		row, err := stmt.Query(id)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer row.Close()

		// Iterate over the rows and scan them into KeyPairs structs
		for row.Next() {
			if err := row.Scan(&subscriber.ID, &subscriber.UserID, &subscriber.ClientID, &subscriber.Key, &subscriber.CreatedOn); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, subscriber)

	}
}

// AddExample handles POST /api/example - CREATE
func AddKeyPair(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		userId := c.Param("id")

		clientID := generateUniqueID()
		key := generateUniqueID()

		// Create a KeyPairs instance with generated IDs and user ID
		subscriber := models.ApiSubscriber{
			UserID:   userId,
			ClientID: clientID,
			Key:      key,
		}

		//query to insert the record into the database
		query := "INSERT INTO api_subscribers (user_id, client_id, key) VALUES ($1, $2, $3)"

		//prepare statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing the query"})
			return
		}
		//close the statement when the surrounding function returns(handler function)
		defer stmt.Close()

		//execute the statement
		_, err = stmt.Exec(subscriber.UserID, subscriber.ClientID, subscriber.Key)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting into the database"})
			return
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, subscriber)

	}
}

func UpdateKeyPair(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		userId := c.Param("id")

		clientID := generateUniqueID()
		key := generateUniqueID()

		// Create a KeyPairs instance with generated IDs and user ID
		subscriber := models.ApiSubscriber{
			UserID:   userId,
			ClientID: clientID,
			Key:      key,
		}

		//query to update the record in the database
		query := "UPDATE api_subscribers SET client_id = $1, key = $2 WHERE user_id=$3"

		//prepare statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing the query"})
			return
		}
		//close the statement when the surrounding function returns(handler function)
		defer stmt.Close()

		//execute the statement
		_, err = stmt.Exec(subscriber.ClientID, subscriber.Key, userId)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in data Update"})
			return
		}

		// return statement
		c.JSON(http.StatusOK, subscriber)
	}
}

func DeleteKeyPair(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		//query to delete the record from the database
		query := "DELETE FROM api_subscribers WHERE user_id = $1"

		//prepare statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing the query"})
			return
		}
		//close the statement when the surrounding function returns(handler function)
		defer stmt.Close()

		//execute the statement
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete examples"})
			return
		}
		// return statement
		c.JSON(http.StatusOK, gin.H{"message": "Example deleted successfully"})

	}
}

// Function to generate a unique ID
func generateUniqueID() string {
	// Generate a new UUID (version 4)
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error generating UUID:", err)
		return "" // Return an empty string in case of error
	}

	// Convert UUID to string
	return id.String()
}
