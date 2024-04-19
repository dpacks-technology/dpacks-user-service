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

// AddBillingProfile POST /api/billing/profile - CREATE
func AddBillingProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userid, _ := c.Get("auth_userId")
		// get the JSON data
		var transaction models.TransactionsModel
		if err := c.ShouldBindJSON(&transaction); err != nil {
			fmt.Printf("%s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the transaction data
		if err := validators.ValidateNames(transaction, true); err != nil {
			fmt.Printf("%s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to insert the transaction
		query := "INSERT INTO billing_profile ( user_id,company_name,street_no,city, postal_code, country, email, payment_method, given_name , last_name, month, year, cvv, terms, transaction_date,card_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13, $14, $15,$16)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(userid, transaction.CompanyName, transaction.StreetNo, transaction.City, transaction.PostalCode,
			transaction.Country, transaction.Email, transaction.PaymentMethod, transaction.GivenName, transaction.LastName, transaction.Month, transaction.Year, transaction.CVV, transaction.Terms, transaction.TransactionDate, transaction.CardNumber)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Plan added successfully"})

	}
}

// GetBillingProfiles handles GET /api/billing/profile/ - READ
func GetBillingProfiles(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM billing_profile ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT * FROM billing_profile WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "status":
				query = "SELECT * FROM billing_profile WHERE company_name LIKE $3 ORDER BY CASE WHEN status = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
		var transaction []models.TransactionsModel

		for rows.Next() {
			var Transactions models.TransactionsModel
			if err := rows.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			transaction = append(transaction, Transactions)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, transaction)

	}
}

// GetBillingProfileById handles GET /api/billing/profile/:id - READ
func GetBillingProfileById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM billing_profile WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var Transactions models.TransactionsModel

		// Scan the row data into the WebpageModel
		err := row.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, Transactions)

	}
}

// GetBillingProfileByStatusCount handles GET /api/billing/profile/status/:status/count - READ
func GetBillingProfileByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM billing_profile"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM billing_profile WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM billing_profile WHERE status IN ($1)"
			args = append(args, 0)

		case "2":
			query = "SELECT COUNT(*) FROM billing_profile WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM billing_profile WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "status":
				query = "SELECT * FROM billing_profile WHERE company_name LIKE $3 ORDER BY CASE WHEN status = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

// GetTransactionByStatus handles GET /api/billing/profile/status/:status - READ
func GetBillingProfileByStatus(db *sql.DB) gin.HandlerFunc {
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

		query = "SELECT * FROM billing_profile ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM billing_profile WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM billing_profile WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)

		case "2":
			query = "SELECT * FROM billing_profile WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM billing_profile WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM billing_profile WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "status":
				query = "SELECT * FROM billing_profile WHERE company_name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
				//case "path":
				//	query = "SELECT * FROM webpages WHERE path LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				//	args = append(args, escapedVal)
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
		var Transaction []models.TransactionsModel

		for rows.Next() {
			var Transactions models.TransactionsModel
			if err := rows.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			Transaction = append(Transaction, Transactions)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, Transaction)

	}
}

// GetBillingProfileDateTime handles GET /api/billing/profile/datetime/:count/:page - READ
func GetBillingProfileDateTime(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM billing_profile ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM billing_profile WHERE date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM billing_profile WHERE id = $5 AND date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "status":
				query = "SELECT * FROM billing_profile WHERE status LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
		var Transaction []models.TransactionsModel

		for rows.Next() {
			var Transactions models.TransactionsModel
			if err := rows.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			Transaction = append(Transaction, Transactions)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all webpages as JSON
		c.JSON(http.StatusOK, Transaction)

	}
}

// GetBillingProfileByDatetimeCount handles GET /api/billing/proifle/datetime/count - READ
func GetBillingProfileByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM billing_profile"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM billing_profile WHERE date_created BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM billing_profile WHERE id = $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, val)
			case "status":
				query = "SELECT COUNT(*) FROM billing_profile WHERE status LIKE $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, escapedVal)
				//case "path":
				//	query = "SELECT COUNT(*) FROM transactions WHERE path LIKE $3 AND date_created BETWEEN $1 AND $2"
				//	args = append(args, escapedVal)
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

func GetBillingProfileCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM billing_profile"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM billing_profile WHERE id = $1"
				args = append(args, val)
			case "status":
				query = "SELECT COUNT(*) FROM billing_profile WHERE status LIKE $1"
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

// EditBillingProfile handles PUT /api/billing/proifle/:id - UPDATE
func EditBillingProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var transaction models.TransactionsModel
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateNames(transaction, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//print the model
		fmt.Printf("%s", transaction)

		// Update the billing details
		_, err := db.Exec("UPDATE billing_profile SET company_name = $1, street_no= $2, city=$3, postal_code= $4, country=$5, email =$6, payment_method=$7, given_name =$8, last_name = $9, month = $10, year = $11, cvv = $12, terms= $13, card_number = $14 WHERE id = $15", transaction.CompanyName, transaction.StreetNo, transaction.City, transaction.PostalCode, transaction.Country, transaction.Email, transaction.PaymentMethod, transaction.GivenName, transaction.LastName, transaction.Month, transaction.Year, transaction.CVV, transaction.Terms, transaction.CardNumber, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

	}
}

// DeleteBillingProfileByID handles DELETE /api/billing/proifle/:id - DELETE
func DeleteBillingProfileByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM billing_profile WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})

	}
}

// DeleteBillingProfileByIDBulk handles DELETE /api/billing/proifle/bulk/:id - DELETE
func DeleteBillingProfileByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM billing_profile WHERE id = $1"

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
		c.JSON(http.StatusOK, gin.H{"message": "Webpage bulk deleted successfully"})

	}
}

// UpdateBillingProfileBulk handles PUT /api/billing/proifle/status/bulk/:id - UPDATE
func UpdateBillingProfileStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var Transaction models.TransactionsModel
		if err := c.ShouldBindJSON(&Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the transaction status in the database
		for _, id := range ids {

			query := "UPDATE billing_profile SET status = $1 WHERE id = $2"

			// Prepare the statement
			stmt, err := db.Prepare(query)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			// Execute the prepared statement with bound parameters
			_, err = stmt.Exec(Transaction.Status, id)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage status updated successfully"})

	}
}

// UpdateBillingProfileStatus handles PUT /api/billing/proifle/status/:id - UPDATE
func UpdateBillingProfileStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var Transaction models.TransactionsModel
		if err := c.ShouldBindJSON(&Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE billing_profile SET status = $1 WHERE id = $2"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(Transaction.Status, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})

	}
}
