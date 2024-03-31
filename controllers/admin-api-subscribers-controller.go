package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddSubscribers handles POST /api/web/webpages - CREATE
func AddSubscribers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var keypair models.KeyPairs
		if err := c.ShouldBindJSON(&keypair); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//// Validate the webpage data
		//if err := validators.ValidateName(keypair, true); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}
		//keypair.ClientID =

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

// GetApiSubscribers handles GET /api/web/pages/ - READ
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
		query := "SELECT * FROM keypairs ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM keypairs WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "user_id":
				query = "SELECT * FROM keypairs WHERE user_id LIKE $3 ORDER BY CASE WHEN user_id = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "client_id":
				query = "SELECT * FROM keypairs WHERE client_id LIKE $3 ORDER BY CASE WHEN client_id = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Return all webpages as JSON
		c.JSON(http.StatusOK, keypairs)

	}
}

// GetWebPageById handles GET /api/web/webpages/:id - READ
func GetApiSubscriberById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM keypairs WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var keypair models.KeyPairs

		// Scan the row data into the WebpageModel
		err := row.Scan(&keypair.ID, &keypair.UserID, &keypair.ClientID, &keypair.Key, &keypair.CreatedOn)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, keypair)

	}
}

//// GetWebPagesByStatusCount handles GET /api/web/webpages/status/:status/count - READ
//func GetWebPagesByStatusCount(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get status parameter (array)
//		statuses := c.Query("status")
//
//		// get query parameters
//		key := c.Query("key")
//		val := c.Query("val")
//
//		var args []interface{}
//		var query string
//
//		query = "SELECT COUNT(*) FROM webpages"
//
//		switch statuses {
//		case "1":
//			query = "SELECT COUNT(*) FROM webpages WHERE status IN ($1)"
//			args = append(args, 1)
//		case "0":
//			query = "SELECT COUNT(*) FROM webpages WHERE status IN ($1)"
//			args = append(args, 0)
//		}
//
//		if val != "" && key != "" {
//
//			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
//
//			switch key {
//			case "id":
//				query = "SELECT COUNT(*) FROM webpages WHERE id = $2 AND status IN ($1)"
//				args = append(args, val)
//			case "name":
//				query = "SELECT COUNT(*) FROM webpages WHERE name LIKE $2 AND status IN ($1)"
//				args = append(args, escapedVal)
//			case "path":
//				query = "SELECT COUNT(*) FROM webpages WHERE path LIKE $2 AND status IN ($1)"
//				args = append(args, escapedVal)
//			}
//		}
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Execute the prepared statement with bound parameters
//		var count int
//		err = stmt.QueryRow(args...).Scan(&count)
//
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Close the statement
//		defer stmt.Close()
//
//		// Return all webpages as JSON
//		c.JSON(http.StatusOK, count)
//
//	}
//}

//
//// GetWebPagesByStatus handles GET /api/web/webpages/status/:status - READ
//func GetWebPagesByStatus(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get status parameter (array)
//		statuses := c.Query("status")
//
//		// get page id parameter
//		page := c.Param("page")
//
//		// get count parameter
//		count := c.Param("count")
//
//		// Convert page and count to integers
//		pageInt, err := strconv.Atoi(page)
//		if err != nil {
//			// Handle error
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		countInt, err := strconv.Atoi(count)
//		if err != nil {
//			// Handle error
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Calculate offset
//		offset := (pageInt - 1) * countInt
//
//		// get query parameters
//		key := c.Query("key")
//		val := c.Query("val")
//
//		var args []interface{}
//		var query string
//
//		query = "SELECT * FROM webpages ORDER BY id LIMIT $1 OFFSET $2"
//		args = append(args, countInt, offset)
//
//		switch statuses {
//		case "1":
//			query = "SELECT * FROM webpages WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
//			args = append(args, 1)
//		case "0":
//			query = "SELECT * FROM webpages WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
//			args = append(args, 0)
//		}
//
//		if val != "" && key != "" {
//
//			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
//
//			switch key {
//			case "id":
//				query = "SELECT * FROM webpages WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
//				query = "SELECT * FROM webpages WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
//				args = append(args, val)
//			case "name":
//				query = "SELECT * FROM webpages WHERE name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
//				args = append(args, escapedVal)
//			case "path":
//				query = "SELECT * FROM webpages WHERE path LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
//				args = append(args, escapedVal)
//			}
//		}
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Execute the prepared statement with bound parameters
//		rows, err := stmt.Query(args...)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		//close the rows when the surrounding function returns(handler function)
//		defer rows.Close()
//
//		// Iterate over the rows and scan them into WebpageModel structs
//		var webpages []models.WebpageModel
//
//		for rows.Next() {
//			var webpage models.WebpageModel
//			if err := rows.Scan(&webpage.ID, &webpage.Name, &webpage.WebID, &webpage.Path, &webpage.Status, &webpage.DateCreated); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
//				return
//			}
//			webpages = append(webpages, webpage)
//		}
//
//		//this runs only when loop didn't work
//		if err := rows.Err(); err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
//			return
//		}
//
//		// Return all webpages as JSON
//		c.JSON(http.StatusOK, webpages)
//
//	}
//}
//
//// GetWebPagesByDatetime handles GET /api/web/webpages/datetime/:count/:page - READ

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
		query := "SELECT * FROM keypairs ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM keypairs WHERE created_on BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM keypairs WHERE id = $5 AND created_on BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "user_id":
				query = "SELECT * FROM keypairs WHERE user_id LIKE $5 AND created_on BETWEEN $3 AND $4 ORDER BY CASE WHEN user_id = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "client_id":
				query = "SELECT * FROM keypairs WHERE client_id LIKE $5 AND created_on BETWEEN $3 AND $4 ORDER BY CASE WHEN client_id = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Return all webpages as JSON
		c.JSON(http.StatusOK, keypairs)

	}
}

// GetWebPagesByDatetimeCount handles GET /api/web/webpages/datetime/count - READ
func GetApiSubscribersByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM keypairs"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM keypairs WHERE created_on BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM keypairs WHERE id = $3 AND created_on BETWEEN $1 AND $2"
				args = append(args, val)
			case "user_id":
				query = "SELECT COUNT(*) FROM keypairs WHERE user_id LIKE $3 AND created_on BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "client_id":
				query = "SELECT COUNT(*) FROM keypairs WHERE client_id LIKE $3 AND created_on BETWEEN $1 AND $2"
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

func GetApiSubscribersCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM keypairs"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM keypairs WHERE id = $1"
				args = append(args, val)
			case "client_id":
				query = "SELECT COUNT(*) FROM keypairs WHERE client_id LIKE $1"
				args = append(args, escapedVal)
			case "user_id":
				query = "SELECT COUNT(*) FROM keypairs WHERE user_id LIKE $1"
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

//// EditWebPage handles PUT /api/web/webpages/:id - UPDATE
//func EditWebPage(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// get the JSON data - only the name
//		var webpage models.WebpageModel
//		if err := c.ShouldBindJSON(&webpage); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		// Validate the webpage data
//		if err := validators.ValidateName(webpage, false); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		// Update the webpage in the database
//		_, err := db.Exec("UPDATE webpages SET name = $1 WHERE id = $2", webpage.Name, id)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})
//
//	}
//}
//

//// DeleteWebPageByID handles DELETE /api/web/webpages/:id - DELETE
//func DeleteWebPageByID(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// query to delete the webpage
//		query := "DELETE FROM webpages WHERE id = $1"
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Execute the prepared statement with bound parameters
//		_, err = stmt.Exec(id)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage deleted successfully"})
//
//	}
//}
//

//// DeleteWebPageByIDBulk handles DELETE /api/web/webpages/bulk/:id - DELETE
//func DeleteWebPageByIDBulk(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get ids array as a parameter as integer
//		id := c.Param("id")
//
//		// Convert the string of ids to an array of ids
//		ids := strings.Split(id, ",")
//
//		// Delete the webpage from the database
//		for _, id := range ids {
//			// query to delete the webpage
//			query := "DELETE FROM webpages WHERE id = $1"
//
//			// Prepare the statement
//			stmt, err := db.Prepare(query)
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				return
//			}
//
//			// Execute the prepared statement with bound parameters
//			_, err = stmt.Exec(id)
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				return
//			}
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage bulk deleted successfully"})
//
//	}
//}
//

//// UpdateWebPageStatusBulk handles PUT /api/web/webpages/status/bulk/:id - UPDATE
//func UpdateWebPageStatusBulk(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// Convert the string of ids to an array of ids
//		ids := strings.Split(id, ",")
//
//		// get the JSON data - only the status
//		var webpage models.WebpageModel
//		if err := c.ShouldBindJSON(&webpage); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		// Update the webpage status in the database
//		for _, id := range ids {
//
//			query := "UPDATE webpages SET status = $1 WHERE id = $2"
//
//			// Prepare the statement
//			stmt, err := db.Prepare(query)
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				return
//			}
//
//			// Execute the prepared statement with bound parameters
//			_, err = stmt.Exec(webpage.Status, id)
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				return
//			}
//
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})
//
//	}
//}
//

//// UpdateWebPageStatus handles PUT /api/web/webpages/status/:id - UPDATE
//func UpdateWebPageStatus(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// get the JSON data - only the status
//		var webpage models.WebpageModel
//		if err := c.ShouldBindJSON(&webpage); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//
//		// query to update the webpage status
//		query := "UPDATE webpages SET status = $1 WHERE id = $2"
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Execute the prepared statement with bound parameters
//		_, err = stmt.Exec(webpage.Status, id)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			return
//		}
//
//		// Return a success message
//		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})
//
//	}
//}
