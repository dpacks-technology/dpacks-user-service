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
		query := "SELECT * FROM keypairs"

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
		var keypairs []models.KeyPairs

		for rows.Next() {
			var keypair models.KeyPairs
			if err := rows.Scan(&keypair.ID, &keypair.UserID, &keypair.ClientID, &keypair.Key, &keypair.CreatedOn); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			keypairs = append(keypairs, keypair)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all keypairs as JSON
		c.JSON(http.StatusOK, keypairs)

	}
}

func GetKeyPairsID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		// Create an empty ExampleModel struct
		var keypair models.KeyPairs

		//query the database for the record with the given ID
		query := "SELECT * FROM keypairs WHERE user_id = $1"

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
			if err := row.Scan(&keypair.ID, &keypair.UserID, &keypair.ClientID, &keypair.Key, &keypair.CreatedOn); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, keypair)

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
		keypair := models.KeyPairs{
			UserID:   userId,
			ClientID: clientID,
			Key:      key,
		}

		// Insert the record into the database
		//_, err := db.Exec("INSERT INTO keypairs (user_id, client_id, key) VALUES ($1, $2, $3)", keypair.UserID, keypair.ClientID, keypair.Key)
		//if err != nil {
		//	fmt.Printf("%s\n", err)
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting into the database"})
		//	return
		//}

		//query to insert the record into the database
		query := "INSERT INTO keypairs (user_id, client_id, key) VALUES ($1, $2, $3)"

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
		_, err = stmt.Exec(keypair.UserID, keypair.ClientID, keypair.Key)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting into the database"})
			return
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, keypair)

	}
}

func UpdateKeyPair(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		userId := c.Param("id")

		clientID := generateUniqueID()
		key := generateUniqueID()

		// Create a KeyPairs instance with generated IDs and user ID
		keypair := models.KeyPairs{
			UserID:   userId,
			ClientID: clientID,
			Key:      key,
		}

		//_, err := db.Exec("UPDATE keypairs SET client_id = $1, key = $2 WHERE user_id=$3", keypair.ClientID, keypair.Key, userId)
		//if err != nil {
		//	fmt.Printf("%s\n", err)
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in data Update"})
		//	return
		//}

		//query to update the record in the database
		query := "UPDATE keypairs SET client_id = $1, key = $2 WHERE user_id=$3"

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
		_, err = stmt.Exec(keypair.ClientID, keypair.Key, userId)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in data Update"})
			return
		}

		// return statement
		c.JSON(http.StatusOK, keypair)
	}
}

func DeleteKeyPair(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		//result, err := db.Exec("DELETE FROM keypairs WHERE user_id = $1", id)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete examples"})
		//	return
		//}
		//
		//rowCount, err := result.RowsAffected()
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//fmt.Printf("Deleted %d rows\n", rowCount)

		//query to delete the record from the database
		query := "DELETE FROM keypairs WHERE user_id = $1"

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
