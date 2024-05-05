package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"dpacks-go-services-template/utils"
	"dpacks-go-services-template/validators"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AddTemplate handles POST /api/marketplace/template - CREATE
func AddTemplate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userid, _ := c.Get("auth_userId")

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

		submittedDate := time.Now()

		mainFileLink := "https://storage.googleapis.com/dpacks-templates.appspot.com/" + template.MainFile

		// query to insert the template
		query := "INSERT INTO templates (name, description, category, mainfile, thmbnlfile, userid, dmessage, price, submitteddate, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(template.Name, template.Description, template.Category, mainFileLink, template.ThmbnlFile, userid, template.DevpDescription, template.Price, submittedDate, 0)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Template submitted successfully"})

	}
}

// GetTemplates handles GET /api/templates/:count/:page - READ
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

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all templates as JSON
		c.JSON(http.StatusOK, templates)

	}
}

// GetTemplatesCount handles GET /api/templates/count - READ
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

		// Return all templates as JSON
		c.JSON(http.StatusOK, count)

	}
}

// GetTemplatesByStatusCount handles GET /api/templates/status/count - READ
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
				query = "SELECT COUNT(*) FROM templates WHERE category LIKE $2 AND status IN ($1)"
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

// GetTemplatesByStatus handles GET /api/templates/status/:count/:page - READ
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
				query = "SELECT * FROM templates WHERE category LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
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

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all templates as JSON
		c.JSON(http.StatusOK, templates)

	}
}

// GetTemplatesByDatetime handles GET /api/templates/datetime/:count/:page - READ
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
			query = "SELECT * FROM templates WHERE submitteddate BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM templates WHERE id = $5 AND submitteddate BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "name":
				query = "SELECT * FROM templates WHERE name LIKE $5 AND submitteddate BETWEEN $3 AND $4 ORDER BY CASE WHEN name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT * FROM templates WHERE category LIKE $5 AND submitteddate BETWEEN $3 AND $4 ORDER BY CASE WHEN category = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all templates as JSON
		c.JSON(http.StatusOK, templates)

	}
}

// GetTemplatesByDatetimeCount handles GET /api/web/templates/datetime/count - READ
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
			query = "SELECT COUNT(*) FROM templates WHERE submitteddate BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM templates WHERE id = $3 AND submitteddate BETWEEN $1 AND $2"
				args = append(args, val)
			case "name":
				query = "SELECT COUNT(*) FROM templates WHERE name LIKE $3 AND submitteddate BETWEEN $1 AND $2"
				args = append(args, escapedVal)
			case "category":
				query = "SELECT COUNT(*) FROM templates WHERE category LIKE $3 AND submitteddate BETWEEN $1 AND $2"
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

		// Return all templates as JSON
		c.JSON(http.StatusOK, count)

	}
}

// UpdateTemplatesStatus handles PUT /api/templates/status/:id - UPDATE
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

		// query to update the template status
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

		// Send email notification to the user
		err = utils.SendEmail("ishini.aponso1230@gmail.com", "New template status change.", "The status of your newly added template has been changed by the marketplace admin. Please check your template.", "small")
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email"})
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Template status updated successfully"})

	}
}

// UpdateTemplatesStatusBulk handles PUT /api/templates/status/bulk/:id - UPDATE
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

		// Update the template status in the database
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

// EditTemplatesD handles PUT /api/templates/:id - UPDATE
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

		//Validate the template data
		if err := validators.ValidateTempEdit(template, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the template in the database
		_, err := db.Exec("UPDATE templates SET name = $1, description = $2, category = $3, dmessage = $4  WHERE id = $5", template.Name, template.Description, template.Category, template.DevpDescription, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Template updated successfully"})

	}
}

// DeleteTemplateByID handles DELETE /api/templates/:id - DELETE
func DeleteTemplateByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the template
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

// DeleteTemplateByIDBulk handles DELETE /api/templates/bulk/:id - DELETE
func DeleteTemplateByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the template from the database
		for _, id := range ids {
			// query to delete the template
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

// GetTemplatesById handles GET /api/template/:id - READ
func GetTemplatesById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM templates WHERE id = $1", id)

		// Create a TemplateModel to hold the data
		var template models.TemplateModel

		// Scan the row data into the TemplateModel
		err := row.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the template as JSON
		c.JSON(http.StatusOK, template)

	}
}

// GetTemplatesBydid handles GET /api/templates/user/:count/:page - READ
func GetTemplatesBydid(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, _ := c.Get("auth_userId")
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
		//UserID := 1

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
		rows, err := stmt.Query(userid, countInt, offset)
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

// DownloadById handles GET /api/templat/:id - READ
func DownloadById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT mainfile FROM templates WHERE id = $1", id)

		// Create a TemplateModel to hold the data
		var template models.TemplateModel

		// Scan the row data into the TemplateModel
		err := row.Scan(&template.MainFile)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the templates as JSON
		c.JSON(http.StatusOK, template)

	}
}

// GetAcceptedTemplates handles GET /api/templates/acceptstatus/:count/:page - READ
//func GetAcceptedTemplates(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get page and count parameters
//		page := c.Param("page")
//		count := c.Param("count")
//
//		// Convert page and count to integers
//		pageInt, err := strconv.Atoi(page)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
//			return
//		}
//
//		countInt, err := strconv.Atoi(count)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
//			return
//		}
//
//		// Calculate offset
//		offset := (pageInt - 1) * countInt
//
//		// Query the database for records based on pagination and userid
//		//query := "SELECT * FROM templates WHERE status = 1 ORDER BY status LIMIT $1 OFFSET $2"
//		query := "SELECT * FROM templates WHERE status = 1 ORDER BY id LIMIT $1 OFFSET $2"
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			return
//		}
//		defer stmt.Close()
//
//		// Execute the prepared statement with bound parameters
//		rows, err := stmt.Query(countInt, offset)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			return
//		}
//		defer rows.Close()
//
//		// Iterate over the rows and scan them into TemplateModel structs
//		var templates []models.TemplateModel
//		for rows.Next() {
//			var template models.TemplateModel
//			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
//				return
//			}
//			templates = append(templates, template)
//		}
//
//		// Check for errors during iteration
//		if err := rows.Err(); err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
//			return
//		}
//
//		// Return templates as JSON
//		c.JSON(http.StatusOK, templates)
//	}
//}

// GetAcceptedTemplates handles GET /api/templates/acceptstatus/:count/:page - READ
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
		query := `SELECT templates.id, templates.name, templates.description, templates.category, templates.mainfile, templates.thmbnlfile, templates.userid, templates.dmessage, templates.price, templates.submitteddate, templates.status, COALESCE(ROUND(AVG(template_ratings.rating), 1), 0) as average_rating
    FROM templates
    LEFT JOIN template_ratings ON templates.id = template_ratings.id
    WHERE templates.status = 1
    GROUP BY templates.id
    ORDER BY templates.id
    LIMIT $1 OFFSET $2`

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
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status, &template.Rating); err != nil {
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

//func GetTemplatesByCategory(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get page and count parameters
//		page := c.Param("page")
//		count := c.Param("count")
//
//		// Convert page and count to integers
//		pageInt, err := strconv.Atoi(page)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
//			return
//		}
//
//		countInt, err := strconv.Atoi(count)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
//			return
//		}
//
//		// Calculate offset
//		offset := (pageInt - 1) * countInt
//
//		// Get categories parameter
//		categories := c.Query("categories")
//
//		// Query the database for records based on pagination and userid
//		query := `SELECT templates.id, templates.name, templates.description, templates.category, templates.mainfile, templates.thmbnlfile, templates.userid, templates.dmessage, templates.price, templates.submitteddate, templates.status, COALESCE(AVG(template_ratings.rating), 0) as average_rating
//    FROM templates
//    LEFT JOIN template_ratings ON templates.id = template_ratings.id
//    WHERE templates.status = 1`
//
//		// If categories are provided, add a WHERE clause to the query
//		if categories != "" {
//			query += " AND templates.category IN (" + categories + ")"
//		}
//
//		query += ` GROUP BY templates.id
//    ORDER BY templates.id
//    LIMIT $1 OFFSET $2`
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			return
//		}
//		defer stmt.Close()
//
//		// Execute the prepared statement with bound parameters
//		rows, err := stmt.Query(countInt, offset)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			return
//		}
//		defer rows.Close()
//
//		// Iterate over the rows and scan them into TemplateModel structs
//		var templates []models.TemplateModel
//		for rows.Next() {
//			var template models.TemplateModel
//			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status, &template.Rating); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
//				return
//			}
//			templates = append(templates, template)
//		}
//
//		// Check for errors during iteration
//		if err := rows.Err(); err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
//			return
//		}
//
//		// Return templates as JSON
//		c.JSON(http.StatusOK, templates)
//	}
//}

// GetTemplatesByCategory handles GET /api/templates/filter/:count/:page/:category - READ
func GetTemplatesByCategory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get page and count parameters
		page := c.Param("page")
		count := c.Param("count")

		// Get category parameter
		category := c.Param("category")

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

		// Query the database for records based on pagination and category
		query := `SELECT templates.id, templates.name, templates.description, templates.category, templates.mainfile, templates.thmbnlfile, templates.userid, templates.dmessage, templates.price, templates.submitteddate, templates.status, COALESCE(AVG(template_ratings.rating), 0) as average_rating
     FROM templates
     LEFT JOIN template_ratings ON templates.id = template_ratings.id
     WHERE templates.status = 1 AND templates.category = $3
     GROUP BY templates.id
     ORDER BY templates.id
     LIMIT $1 OFFSET $2`

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer stmt.Close()

		// Execute the prepared statement with bound parameters
		rows, err := stmt.Query(countInt, offset, category)
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
			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status, &template.Rating); err != nil {
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

// GetbySearchListingPage handles GET /api/templates/search/:count/:page - READ
//func SearchTemplatesByCategory(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get category parameter
//		category := c.Query("category")
//
//		// Query the database for records based on category
//		query := `SELECT templates.id, templates.name, templates.description, templates.category, templates.mainfile, templates.thmbnlfile, templates.userid, templates.dmessage, templates.price, templates.submitteddate, templates.status, COALESCE(AVG(template_ratings.rating), 0) as average_rating
//      FROM templates
//      LEFT JOIN template_ratings ON templates.id = template_ratings.id
//      WHERE templates.status = 1 AND templates.category = $1
//      GROUP BY templates.id
//      ORDER BY templates.id`
//
//		// Prepare the statement
//		stmt, err := db.Prepare(query)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			return
//		}
//		defer stmt.Close()
//
//		// Execute the prepared statement with bound parameters
//		rows, err := stmt.Query(category)
//		if err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
//			return
//		}
//		defer rows.Close()
//
//		// Iterate over the rows and scan them into TemplateModel structs
//		var templates []models.TemplateModel
//		for rows.Next() {
//			var template models.TemplateModel
//			if err := rows.Scan(&template.Id, &template.Name, &template.Description, &template.Category, &template.MainFile, &template.ThmbnlFile, &template.UserID, &template.DevpDescription, &template.Price, &template.Sdate, &template.Status, &template.Rating); err != nil {
//				fmt.Printf("%s\n", err)
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
//				return
//			}
//			templates = append(templates, template)
//		}
//
//		// Check for errors during iteration
//		if err := rows.Err(); err != nil {
//			fmt.Printf("%s\n", err)
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
//			return
//		}
//
//		// Return templates as JSON
//		c.JSON(http.StatusOK, templates)
//	}
//}

// UploadTemplate handles POST /api/template/upload - CREATE
func UploadTemplate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate random name to the file with a random string and the current timestamp
		var tempName = "template_" + utils.GenerateRandomString(10) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".zip"

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = utils.UploadTemplate(tempName, file)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Save the file to the server
		if err := c.SaveUploadedFile(file, "./uploads/"+tempName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the file"})
			return
		}

		fmt.Printf("File %s uploaded successfully\n", tempName)
		// Return a success response
		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"fileName": tempName,
		})
	}
}

// UploadThumbImg handles POST /api/template/image/upload - CREATE
func UploadThumbImg(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate random name to the file with a random string and the current timestamp
		var thumbName = "templateThumb_" + utils.GenerateRandomString(10) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".png"

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = utils.UploadThumbImg(thumbName, file)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Save the file to the server
		if err := c.SaveUploadedFile(file, "./uploads/"+thumbName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the file"})
			return
		}

		fmt.Printf("File %s uploaded successfully\n", thumbName)
		// Return a success response
		c.JSON(http.StatusOK, gin.H{
			"message":  "File uploaded successfully",
			"fileName": thumbName,
		})
	}
}

// AddRatings handles POST /api/template/rating - CREATE
func AddRatings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userid, _ := c.Get("auth_userId")

		// get the JSON data
		var templateR models.TemplateRatingsModel
		if err := c.ShouldBindJSON(&templateR); err != nil {
			fmt.Printf("%s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Validate the data
		//if err := validators.ValidateTemp(template, true); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}

		ratedDate := time.Now()

		// query to insert the template
		query := "INSERT INTO template_ratings (id, user_id, rating, rating_date) VALUES ($1, $2, $3, $4)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(templateR.TID, userid, templateR.Rating, ratedDate)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Rate added successfully"})

	}
}
