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

// AddAdminUser - handles POST - CREATE
func AddAdminUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var admin models.AdminUserModel
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the admin data
		if err := validators.ValidateAdmin(admin, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the admin
		query := "INSERT INTO admin_user (name, phone, email, password) VALUES ($1, $2, $3, $4)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(admin.Name, admin.Phone, admin.Email, admin.Password)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Admin added successfully"})

	}
}

// GetAdmins - handles GET  - READ
func GetAdmins(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM admin_user ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM admin_user WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM admin_user WHERE name LIKE $3 ORDER BY CASE WHEN name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "phone":
				query = "SELECT * FROM admin_user WHERE phone LIKE $3 ORDER BY CASE WHEN phone = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT * FROM admin_user WHERE email LIKE $3 ORDER BY CASE WHEN email = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into adminModel structs
		var admins []models.AdminUserModel

		for rows.Next() {
			var admin models.AdminUserModel
			if err := rows.Scan(&admin.ID, &admin.Name, &admin.Phone, &admin.Email, &admin.Password, &admin.AddedOn, &admin.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			admins = append(admins, admin)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all admins as JSON
		c.JSON(http.StatusOK, admins)

	}
}

// GetAdminById - handles GET  - READ
func GetAdminById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM admin_user WHERE id = $1", id)

		// Create a AdminUserModel to hold the data
		var admin models.AdminUserModel

		// Scan the row data into the AdminUserModel
		err := row.Scan(&admin.ID, &admin.Name, &admin.Phone, &admin.Email, &admin.Password, &admin.AddedOn, &admin.Status)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the admin as JSON
		c.JSON(http.StatusOK, admin)

	}
}

// GetAdminsByStatusCount - handles GET - READ
func GetAdminsByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM admin_user"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM admin_user WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM admin_user WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM admin_user WHERE id = $2 AND status IN ($1)"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM admin_user WHERE name LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "phone":
				query = "SELECT COUNT(*) FROM admin_user WHERE phone LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT COUNT(*) FROM admin_user WHERE email LIKE $2 AND status IN ($1)"
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

		// Return all admins as JSON
		c.JSON(http.StatusOK, count)

	}
}

// GetAdminsByStatus handles GET - READ
func GetAdminsByStatus(db *sql.DB) gin.HandlerFunc {
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

		query = "SELECT * FROM admin_user ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM admin_user WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM admin_user WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM admin_user WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM admin_user WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM admin_user WHERE name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "phone":
				query = "SELECT * FROM admin_user WHERE phone LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT * FROM admin_user WHERE email LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into AdminUserModel structs
		var admins []models.AdminUserModel

		for rows.Next() {
			var admin models.AdminUserModel
			if err := rows.Scan(&admin.ID, &admin.Name, &admin.Phone, &admin.Email, &admin.Password, &admin.AddedOn, &admin.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			admins = append(admins, admin)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all admins as JSON
		c.JSON(http.StatusOK, admins)

	}
}

// GetAdminsByDatetime handles GET - READ
func GetAdminsByDatetime(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM admin_user ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM admin_user WHERE added_on BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM admin_user WHERE id = $5 AND added_on BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM admin_user WHERE name LIKE $5 AND added_on BETWEEN $3 AND $4 ORDER BY CASE WHEN name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "phone":
				query = "SELECT * FROM admin_user WHERE phone LIKE $5 AND added_on BETWEEN $3 AND $4 ORDER BY CASE WHEN phone = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT * FROM admin_user WHERE email LIKE $5 AND added_on BETWEEN $3 AND $4 ORDER BY CASE WHEN email = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		// Iterate over the rows and scan them into AdminUserModel structs
		var admins []models.AdminUserModel

		for rows.Next() {
			var admin models.AdminUserModel
			if err := rows.Scan(&admin.ID, &admin.Name, &admin.Phone, &admin.Email, &admin.Password, &admin.AddedOn, &admin.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			admins = append(admins, admin)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all admins as JSON
		c.JSON(http.StatusOK, admins)

	}
}

// GetAdminsByDatetimeCount handles GET - READ
func GetAdminsByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM admin_user"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM admin_user WHERE added_on BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM admin_user WHERE id = $3 AND added_on BETWEEN $1 AND $2"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM admin_user WHERE name LIKE $3 AND added_on BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "phone":
				query = "SELECT COUNT(*) FROM admin_user WHERE phone LIKE $3 AND added_on BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT COUNT(*) FROM admin_user WHERE email LIKE $3 AND added_on BETWEEN $1 AND $2"
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

		// Return all admins as JSON
		c.JSON(http.StatusOK, count)

	}
}

func GetAdminsCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM admin_user"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM admin_user WHERE id = $1"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM admin_user WHERE name LIKE $1"
				args = append(args, escapedVal)
			case "phone":
				query = "SELECT COUNT(*) FROM admin_user WHERE phone LIKE $1"
				args = append(args, escapedVal)
			case "email":
				query = "SELECT COUNT(*) FROM admin_user WHERE email LIKE $1"
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

		// Return all admins as JSON
		c.JSON(http.StatusOK, count)

	}
}

// EditAdmin - handles PUT  - UPDATE
func EditAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var admin models.AdminUserModel
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the admin data
		if err := validators.ValidateAdmin(admin, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the admin in the database
		_, err := db.Exec("UPDATE admin_user SET name = $1, phone = $2, email = $3, password = $4 WHERE id = $5", admin.Name, admin.Phone, admin.Email, admin.Password, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})

	}
}

// DeleteAdminByID - handles - DELETE
func DeleteAdminByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the admin
		query := "DELETE FROM admin_user WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})

	}
}

// DeleteAdminByIDBulk handles DELETE - DELETE
func DeleteAdminByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the admin from the database
		for _, id := range ids {
			// query to delete the admin
			query := "DELETE FROM admin_user WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Admin bulk deleted successfully"})

	}
}

// UpdateAdminStatusBulk  handles PUT  - UPDATE
func UpdateAdminStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var admin models.AdminUserModel
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the admin status in the database
		for _, id := range ids {

			query := "UPDATE admin_user SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(admin.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Admin status updated successfully"})

	}
}

// UpdateAdminStatus handles PUT - UPDATE
func UpdateAdminStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var admin models.AdminUserModel
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the admin status
		query := "UPDATE admin_user SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(admin.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Admin status updated successfully"})

	}
}
