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

// AddTemplate handles POST /api/marketplace/template - CREATE
func AddTemplate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var template models.TemplateModel
		if err := c.ShouldBindJSON(&template); err != nil {
			fmt.Printf("%s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Validate the data
		if err := validators.ValidateTemp(template, true); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the template
		query := "INSERT INTO templates (name, description, category, mainfile, thmbnlfile, dmessage, price, submitteddate, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(template.Name, template.Description, template.Category, template.MainFile, template.ThmbnlFile, template.DevpDescription, template.Price, template.Sdate, 0)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Template submitted successfully"})

	}
}

// GetTemplates handles GET /api/web/webpages/status/:status/count - READ
func GetTemplates(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM templates ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM templates WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM templates WHERE name LIKE $3 ORDER BY CASE WHEN name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT * FROM templates WHERE category LIKE $3 ORDER BY CASE WHEN category = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
		var templates []models.TemplateModel

		for rows.Next() {
			var template models.TemplateModel
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, templates)

	}
}

func GetTemplatesCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM templates"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM templates WHERE id = $1"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM templates WHERE name LIKE $1"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT COUNT(*) FROM templates WHERE path LIKE $1"
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

// GetWebPagesByStatusCount handles GET /api/web/webpages/status/:status/count - READ
func GetTemplatesByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM templates"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM templates WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM templates WHERE status IN ($1)"
			args = append(args, 0)
		case "2":
			query = "SELECT COUNT(*) FROM templates WHERE status IN ($1)"
			args = append(args, 2)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM templates WHERE id = $2 AND status IN ($1)"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM templates WHERE name LIKE $2 AND status IN ($1)"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT COUNT(*) FROM templates WHERE path LIKE $2 AND status IN ($1)"
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
func GetTemplatesByStatus(db *sql.DB) gin.HandlerFunc {
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

		query = "SELECT * FROM templates ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM templates WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM templates WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		case "2":
			query = "SELECT * FROM templates WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 2)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM templates WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM templates WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM templates WHERE name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT * FROM templates WHERE path LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
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
		var templates []models.TemplateModel

		for rows.Next() {
			var template models.TemplateModel
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, templates)

	}
}

func GetTemplatesByDatetime(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM templates ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM templates WHERE date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM templates WHERE id = $5 AND date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM templates WHERE name LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT * FROM templates WHERE path LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN path = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
		var templates []models.TemplateModel

		for rows.Next() {
			var template models.TemplateModel
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, templates)

	}
}

// GetWebPagesByDatetimeCount handles GET /api/web/webpages/datetime/count - READ
func GetTemplatesByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM templates"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM templates WHERE date_created BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM templates WHERE id = $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM templates WHERE name LIKE $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT COUNT(*) FROM templates WHERE path LIKE $3 AND date_created BETWEEN $1 AND $2"
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

// UpdateWebPageStatus handles PUT /api/web/webpages/status/:id - UPDATE
func UpdateTemplatesStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var template models.TemplateModel
		if err := c.ShouldBindJSON(&template); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE templates SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(template.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Template status updated successfully"})

	}
}

// UpdateWebPageStatusBulk handles PUT /api/web/webpages/status/bulk/:id - UPDATE
func UpdateTemplatesStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var template models.TemplateModel
		if err := c.ShouldBindJSON(&template); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage status in the database
		for _, id := range ids {

			query := "UPDATE templates SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(template.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Template status updated successfully"})

	}
}

func EditTemplatesD(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var template models.TemplateModel
		if err := c.ShouldBindJSON(&template); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		//if err := validators.ValidateTemp(template, false); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}

		// Update the webpage in the database
		_, err := db.Exec("UPDATE templates SET description = $1, category = $2, dmessage = $3  WHERE id = $4", template.Description, template.Category, template.DevpDescription, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Template updated successfully"})

	}
}

// DeleteWebPageByID handles DELETE /api/web/webpages/:id - DELETE
func DeleteTemplateByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM templates WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})

	}
}

// DeleteWebPageByIDBulk handles DELETE /api/web/webpages/bulk/:id - DELETE
func DeleteTemplateByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM templates WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Template bulk deleted successfully"})

	}
}

// GetWebPageById handles GET /api/web/webpages/:id - READ
func GetTemplatesById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM templates WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var template models.TemplateModel

		// Scan the row data into the WebpageModel
		err := row.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, template)

	}
}

func GetTemplatesBydid(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get page and count parameters
		page := c.Param("page")
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// Get query parameters
		//userid := c.Query("userid")
		UserID := 1

		// Query the database for records based on pagination and userid
		query := "SELECT * FROM templates WHERE userid = $1 ORDER BY userid LIMIT $2 OFFSET $3"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(UserID, countInt, offset)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		// Iterate over the rows and scan them into TemplateModel structs
		var templates []models.TemplateModel
		for rows.Next() {
			var template models.TemplateModel
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		// Check for errors during iteration
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return templates as JSON
		c.JSON(http.StatusOK, templates)
	}
}

func DownloadById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT mainfile FROM templates WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var template models.TemplateModel

		// Scan the row data into the WebpageModel
		err := row.Scan(&template.MainFile)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, template)

	}
}

func GetAcceptedTemplates(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get page and count parameters
		page := c.Param("page")
		count := c.Param("count")

		// Convert page and count to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}

		countInt, err := strconv.Atoi(count)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}

		// Calculate offset
		offset := (pageInt - 1) * countInt

		// Query the database for records based on pagination and userid
		query := "SELECT * FROM templates WHERE status = 1 ORDER BY status LIMIT $1 OFFSET $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(countInt, offset)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		// Iterate over the rows and scan them into TemplateModel structs
		var templates []models.TemplateModel
		for rows.Next() {
			var template models.TemplateModel
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		// Check for errors during iteration
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return templates as JSON
		c.JSON(http.StatusOK, templates)
	}
}

// AddWebPage handles POST /api/web/webpages - CREATE
func AddTemplateRatings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the JSON data
		var template models.TemplateRatingsModel
		if err := c.ShouldBindJSON(&template); err != nil {
			fmt.Printf("%s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Validate the webpage data
		//if err := validators.ValidateTemp(template, true); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}

		// query to insert the webpage
		query := "INSERT INTO template_ratings (id, rating) VALUES ($1, $2)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(template.TID, template.Rating)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Rating submitted successfully"})

	}
}

//	func GetRatingSumAndCount(db *sql.DB) gin.HandlerFunc {
//		return func(c *gin.Context) {
//			// get query parameters
//			key := c.Query("key")
//			val := c.Query("val")
//			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
//
//			var args []interface{}
//
//			// Query the database for sum and count of ratings
//			query := `SELECT t.id, t.name, t.description, t.category, t.mainfile, t.thmbnlfile, t.userid, t.dmessage, t.price, t.submitteddate, t.status,
//	          COALESCE(SUM(tr.rating), 0) AS total_ratings, COUNT(tr.rating) AS rating_count
//	   FROM templates t
//	   LEFT JOIN template_ratings tr ON t.id = tr.id
//	   GROUP BY t.id, t.name, t.description, t.category, t.mainfile, t.thmbnlfile, t.userid, t.dmessage, t.price, t.submitteddate, t.status
//	   `
//
//			if val != "" && key != "" {
//				switch key {
//				case "id":
//					query += " WHERE t.id = $1"
//					args = append(args, val)
//				case "name":
//					query += " WHERE t.name LIKE $1"
//					args = append(args, escapedVal)
//				case "category":
//					query += " WHERE t.category LIKE $1"
//					args = append(args, escapedVal)
//				}
//			}
//
//			// Prepare the statement
//			stmt, err := db.Prepare(query)
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error preparing statement"})
//				return
//			}
//			defer stmt.Close()
//
//			// Execute the prepared statement with bound parameters
//			rows, err := stmt.Query(args...)
//			if err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing query"})
//				return
//			}
//			defer rows.Close()
//
//			// Define variables to accumulate the sum and count of ratings
//			var totalRatings int
//			var ratingCount int
//
//			// Iterate through the rows and accumulate the sum and count of ratings
//			for rows.Next() {
//				var id int
//				var name, description, category, mainfile, thmbnlfile, dmessage string
//				var price float64
//				var submitteddate int64
//				var status int
//				var userID sql.NullInt64 // Use sql.NullInt64 for the userid column
//				err := rows.Scan(&id, &name, &description, &category, &mainfile, &thmbnlfile, &userID, &dmessage, &price, &submitteddate, &status, &totalRatings, &ratingCount)
//				if err != nil {
//					fmt.Printf("%s\n", err)
//					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
//					return
//				}
//
//				// Convert sql.NullInt64 to int if it's valid, otherwise assign a default value (e.g., 0)
//				var userIDValue int
//				if userID.Valid {
//					userIDValue = int(userID.Int64)
//				} else {
//					userIDValue = 0 // Assign a default value
//				}
//
//				// You can use the retrieved values as needed
//			}
//
//			// Return sum and count of ratings as JSON
//			c.JSON(http.StatusOK, gin.H{"total_ratings": totalRatings, "rating_count": ratingCount})
//		}
//	}
func GetbySearchListingPage(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM templates ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "name":
				query = "SELECT * FROM templates WHERE name LIKE $3 ORDER BY CASE WHEN name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT * FROM templates WHERE category LIKE $3 ORDER BY CASE WHEN category = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
		var templates []models.TemplateModel

		for rows.Next() {
			var template models.TemplateModel
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			templates = append(templates, template)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, templates)

	}
}
