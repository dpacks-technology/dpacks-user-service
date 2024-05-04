package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTotalUserCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM users"

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

		// Return all users as JSON
		c.JSON(http.StatusOK, count)

	}
}

func GetTotalWebsitesCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM sites"

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

		// Return all users as JSON
		c.JSON(http.StatusOK, count)

	}
}

func GetTotalApiSubscribersCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM api_subscribers"

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

		// Return all users as JSON
		c.JSON(http.StatusOK, count)

	}
}

func GetTotalMarketplaceUsersCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(DISTINCT userid) FROM templates"

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

		// Return all users as JSON
		c.JSON(http.StatusOK, count)

	}
}
