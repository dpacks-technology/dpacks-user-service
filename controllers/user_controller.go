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

// GetUsers handles GET - READ
func GetUsers(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM users WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "first_name":
				query = "SELECT * FROM users WHERE first_name LIKE $3 ORDER BY CASE WHEN first_name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "last_name":
				query = "SELECT * FROM users WHERE last_name LIKE $3 ORDER BY CASE WHEN last_name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT * FROM users WHERE email LIKE $3 ORDER BY CASE WHEN email = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into userModel structs
		var users []models.UserModel

		for rows.Next() {
			var user models.UserModel
			if err := rows.Scan(&user.ID, &user.Code, &user.DOB, &user.FirstName, &user.ForgotCode, &user.ForgotCodeExpire, &user.Gender, &user.InitDate, &user.LastName, &user.Password, &user.Phone, &user.Status, &user.UserKey, &user.Email, &user.VerificationExp); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			users = append(users, user)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all users as JSON
		c.JSON(http.StatusOK, users)

	}
}

//// GetUserById handles GET /api/web/webpages/:id - READ
//func GetWebPageById(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		// get id parameter
//		id := c.Param("id")
//
//		// Query the database for a single record
//		row := db.QueryRow("SELECT * FROM users WHERE id = $1", id)
//
//		// Create a UserModel to hold the data
//		var user models.UserModel
//
//		// Scan the row data into the UserModel
//		err := row.Scan(&user.ID, &user.Code, &user.DOB, &user.FirstName, &user.ForgotCode, &user.ForgotCodeExpire, &user.Gender, &user.InitDate, &user.LastName, &user.Password, &user.Phone, &user.Status, &user.UserKey, &user.Email, &user.VerificationExp)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
//			return
//		}
//
//		// Return the user as JSON
//		c.JSON(http.StatusOK, user)
//
//	}
//}

// GetUsersByStatusCount handles GET /api/user/users/status/:status/count - READ
func GetUsersByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM users"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM users WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM users WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM users WHERE id = $2 AND status IN ($1)"
				args = append(args, val)
			case "first_name":
				query = "SELECT COUNT(*) FROM users WHERE first_name LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "last_name":
				query = "SELECT COUNT(*) FROM users WHERE last_name LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT COUNT(*) FROM users WHERE email LIKE $2 AND status IN ($1)"
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

		// Return all users as JSON
		c.JSON(http.StatusOK, count)

	}
}

// GetUsersByStatus handles GET - READ
func GetUsersByStatus(db *sql.DB) gin.HandlerFunc {
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

		query = "SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM users WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM users WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM users WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM users WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "first_name":
				query = "SELECT * FROM users WHERE first_name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "last_name":
				query = "SELECT * FROM users WHERE last_name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT * FROM users WHERE email LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into UserModel structs
		var users []models.UserModel

		for rows.Next() {
			var user models.UserModel
			if err := rows.Scan(&user.ID, &user.Code, &user.DOB, &user.FirstName, &user.ForgotCode, &user.ForgotCodeExpire, &user.Gender, &user.InitDate, &user.LastName, &user.Password, &user.Phone, &user.Status, &user.UserKey, &user.Email, &user.VerificationExp); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			users = append(users, user)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all users as JSON
		c.JSON(http.StatusOK, users)

	}
}

// GetUsersByDatetime handles GET - READ
func GetUsersByDatetime(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM users WHERE init_date BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM users WHERE id = $5 AND init_date BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "first_name":
				query = "SELECT * FROM users WHERE first_name LIKE $5 AND init_date BETWEEN $3 AND $4 ORDER BY CASE WHEN first_name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "last_name":
				query = "SELECT * FROM users WHERE last_name LIKE $5 AND init_date BETWEEN $3 AND $4 ORDER BY CASE WHEN last_name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT * FROM users WHERE email LIKE $5 AND init_date BETWEEN $3 AND $4 ORDER BY CASE WHEN email = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into UsersModel structs
		var users []models.UserModel

		for rows.Next() {
			var user models.UserModel
			if err := rows.Scan(&user.ID, &user.Code, &user.DOB, &user.FirstName, &user.ForgotCode, &user.ForgotCodeExpire, &user.Gender, &user.InitDate, &user.LastName, &user.Password, &user.Phone, &user.Status, &user.UserKey, &user.Email, &user.VerificationExp); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			users = append(users, user)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all users as JSON
		c.JSON(http.StatusOK, users)

	}
}

// GetUsersByDatetimeCount handles GET - READ
func GetUsersByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM users"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM users WHERE init_date BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM users WHERE id = $3 AND init_date BETWEEN $1 AND $2"
				args = append(args, val)
			case "first_name":
				query = "SELECT COUNT(*) FROM users WHERE first_name LIKE $3 AND init_date BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "last_name":
				query = "SELECT COUNT(*) FROM users WHERE last_name LIKE $3 AND init_date BETWEEN $1 AND $2"
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

		// Return all users as JSON
		c.JSON(http.StatusOK, count)

	}
}

func GetUsersCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM users"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM users WHERE id = $1"
				args = append(args, val)
			case "first_name":
				query = "SELECT COUNT(*) FROM users WHERE first_name LIKE $1"
				args = append(args, escapedVal)
			case "last_name":
				query = "SELECT COUNT(*) FROM users WHERE last_name LIKE $1"
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

// DeleteWebPageByID handles DELETE /api/web/webpages/:id - DELETE
func DeleteUserByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM users WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

	}
}

// DeleteUserByIDBulk handles DELETE - DELETE
func DeleteUserByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the user from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM users WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "User bulk deleted successfully"})

	}
}

// UpdateUserStatusBulk handles PUT - UPDATE
func UpdateUserStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var user models.UserModel
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the user status in the database
		for _, id := range ids {

			query := "UPDATE users SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(user.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})

	}
}

// UpdateUserStatus handles PUT - UPDATE
func UpdateUserStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var user models.UserModel
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the user status
		query := "UPDATE users SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(user.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})

	}
}
