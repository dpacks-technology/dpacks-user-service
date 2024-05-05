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

//function to get site specific storage
//func GetSitesStorage(db *sql.DB) gin.HandlerFunc {
//	// Return a handler function
//	return func(c *gin.Context) {
//
//		// Execute the SQL query
//		rows, err := db.Query(`
//            SELECT sites.name AS site_names,
//            (CAST(sum(size) AS FLOAT8)/1048576) AS size_in_mb
//			FROM data_packets
//			INNER JOIN public.sites ON public.sites.id::text = data_packets.site
//			GROUP BY sites.id;`)
//		if err != nil {
//			// Handle any errors
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
//			return
//		}
//		defer rows.Close()
//
//		// Iterate over the result set
//		var results []gin.H
//		for rows.Next() {
//			var siteNames string
//			var site_sum float64
//			if err := rows.Scan(&siteNames, &site_sum); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning database rows"})
//				return
//			}
//
//			results = append(results, gin.H{"site_name": siteNames, "site_sum": site_sum})
//		}
//
//		// Return the results as JSON
//		c.JSON(http.StatusOK, results)
//	}
//}

func GetTotalUsedStorage(db *sql.DB) gin.HandlerFunc {
	// Return a handler function
	return func(c *gin.Context) {

		//query the database for the record with the given ID
		query := "SELECT sum(size) AS total_sum from data_packets"

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
		row, err := stmt.Query()
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer row.Close()

		var total_sum float64
		for row.Next() {
			if err := row.Scan(&total_sum); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
		}

		//total_sum = total_sum / 1048576

		// Return the example as JSON
		c.JSON(http.StatusOK, total_sum)
	}
}
