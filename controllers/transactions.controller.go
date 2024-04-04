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

// AddTransactionhandles POST /api/billing/transaction - CREATE
func AddTransaction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

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
		query := "INSERT INTO transactions (id,user_id,plan_id,plan_name, amount,company_name,street_no,city, postal_code, country, email, payment_method, given_name , last_name, month, year, cvv, terms, transaction_date, status,card_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11,$12,$13, $14, $15, $16, $17,$18,$19,$20,$21)"

		// Prepare the statement
		stmt, err := db.Prepare(query)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Execute the prepared statement with bound parameters
		_, err = stmt.Exec(transaction.TransactionID, transaction.UserID, transaction.PlanID, transaction.PlanName, transaction.Amount, transaction.CompanyName, transaction.StreetNo, transaction.City, transaction.PostalCode,
			transaction.Country, transaction.Email, transaction.PaymentMethod, transaction.GivenName, transaction.LastName, transaction.Month, transaction.Year, transaction.CVV, transaction.Terms, transaction.TransactionDate, transaction.Status, transaction.CardNumber)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusCreated, gin.H{"message": "Plan added successfully"})

	}
}

// GetTransactions handles GET /api/billing/transactions/ - READ
func GetTransactions(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM transactions ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if val != "" && key != "" {
			switch key {
			case "transaction_id":
				query = "SELECT * FROM transactions WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "plan_name":
				query = "SELECT * FROM transactions WHERE plan_name LIKE $3 ORDER BY CASE WHEN plan_name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
			if err := rows.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.PlanID, &Transactions.PlanName, &Transactions.Amount, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber); err != nil {
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

// GetTansactionsById handles GET /api/web/webpages/:id - READ
func GetTansactionsById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Query the database for a single record
		row := db.QueryRow("SELECT * FROM transaction WHERE id = $1", id)

		// Create a WebpageModel to hold the data
		var Transactions models.TransactionsModel

		// Scan the row data into the WebpageModel
		err := row.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.PlanID, &Transactions.PlanName, &Transactions.Amount, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row from the database"})
			return
		}

		// Return the webpage as JSON
		c.JSON(http.StatusOK, Transactions)

	}
}

// GetTransactionByStatusCount handles GET /api/web/webpages/status/:status/count - READ
func GetTransactionByStatusCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get status parameter (array)
		statuses := c.Query("status")

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM transactions"

		switch statuses {
		case "1":
			query = "SELECT COUNT(*) FROM transactions WHERE status IN ($1)"
			args = append(args, 1)
		case "0":
			query = "SELECT COUNT(*) FROM transactions WHERE status IN ($1)"
			args = append(args, 0)

		case "2":
			query = "SELECT COUNT(*) FROM transactions WHERE status IN ($1)"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM transactions WHERE id = $3 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "plan_name":
				query = "SELECT * FROM transactions WHERE name LIKE $3 ORDER BY CASE WHEN name = $3 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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

// GetTransactionByStatus handles GET /api/web/webpages/status/:status - READ
func GetTransactionByStatus(db *sql.DB) gin.HandlerFunc {
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

		query = "SELECT * FROM transactions ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		switch statuses {
		case "1":
			query = "SELECT * FROM transactions WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 1)
		case "0":
			query = "SELECT * FROM transactions WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)

		case "2":
			query = "SELECT * FROM transactions WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, 0)
		}

		if val != "" && key != "" {

			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"

			switch key {
			case "id":
				query = "SELECT * FROM transactions WHERE status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				query = "SELECT * FROM transactions WHERE id = $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "plan_name":
				query = "SELECT * FROM transactions WHERE name LIKE $4 AND status IN ($3) ORDER BY id LIMIT $1 OFFSET $2"
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
			if err := rows.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.PlanID, &Transactions.PlanName, &Transactions.Amount, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber); err != nil {
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

// GetTransactionDateTime handles GET /api/web/webpages/datetime/:count/:page - READ
func GetTransactionDateTime(db *sql.DB) gin.HandlerFunc {
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
		query := "SELECT * FROM transactions ORDER BY id LIMIT $1 OFFSET $2"
		args = append(args, countInt, offset)

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT * FROM transactions WHERE date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT * FROM transactions WHERE id = $5 AND date_created BETWEEN $3 AND $4 ORDER BY id LIMIT $1 OFFSET $2"
				args = append(args, val)
			case "plan_name":
				query = "SELECT * FROM transactions WHERE plan_name LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN name = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
				args = append(args, escapedVal)
				//case "path":
				//	query = "SELECT * FROM transactions WHERE path LIKE $5 AND date_created BETWEEN $3 AND $4 ORDER BY CASE WHEN path = $5 THEN 1 ELSE 2 END, id LIMIT $1 OFFSET $2"
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
			if err := rows.Scan(&Transactions.TransactionID, &Transactions.UserID, &Transactions.PlanID, &Transactions.PlanName, &Transactions.Amount, &Transactions.CompanyName, &Transactions.StreetNo, &Transactions.City, &Transactions.PostalCode, &Transactions.Country, &Transactions.Email, &Transactions.PaymentMethod, &Transactions.GivenName, &Transactions.LastName, &Transactions.Month, &Transactions.Year, &Transactions.CVV, &Transactions.Terms, &Transactions.TransactionDate, &Transactions.Status, &Transactions.CardNumber); err != nil {
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

// GetTransactionByDatetimeCount handles GET /api/web/webpages/datetime/count - READ
func GetTransactionByDatetimeCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get query parameters
		start := c.Query("start")
		end := c.Query("end")
		key := c.Query("key")
		val := c.Query("val")

		var args []interface{}
		var query string

		query = "SELECT COUNT(*) FROM transactions"

		if start != "" && end != "" && val != "null" && key != "null" {
			query = "SELECT COUNT(*) FROM transactions WHERE date_created BETWEEN $1 AND $2"
			args = append(args, start, end)
		}

		if val != "" && key != "" {
			escapedVal := "%" + strings.ReplaceAll(val, "_", "\\_") + "%"
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM transactions WHERE id = $3 AND date_created BETWEEN $1 AND $2"
				args = append(args, val)
			case "plan_name":
				query = "SELECT COUNT(*) FROM transactions WHERE plan_name LIKE $3 AND date_created BETWEEN $1 AND $2"
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

func GetTransactionCount(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var count int

		// get query parameters
		key := c.Query("key")
		val := c.Query("val")
		escapedVal := strings.ReplaceAll(val, "_", "\\_") + "%"

		var args []interface{}

		// Query the database for records based on pagination
		query := "SELECT COUNT(*) FROM transactions"

		if val != "" && key != "" {
			switch key {
			case "id":
				query = "SELECT COUNT(*) FROM transactions WHERE id = $1"
				args = append(args, val)
			case "plan_name":
				query = "SELECT COUNT(*) FROM transactions WHERE plan_name LIKE $1"
				args = append(args, escapedVal)
				//case "path":
				//	query = "SELECT COUNT(*) FROM transactions WHERE path LIKE $1"
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

// EditTransaction handles PUT /api/web/webpages/:id - UPDATE
func EditTransaction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the name
		var webpage models.WebpageModel
		if err := c.ShouldBindJSON(&webpage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the webpage data
		if err := validators.ValidateName(webpage, false); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the billing details
		_, err := db.Exec("UPDATE transactions SET name = $1 WHERE id = $2", webpage.Name, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		// Return a success message
		c.JSON(http.StatusOK, gin.H{"message": "Webpage updated successfully"})

	}
}

// DeleteTransactionByID handles DELETE /api/web/webpages/:id - DELETE
func DeleteTransactionByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// query to delete the webpage
		query := "DELETE FROM transactions WHERE id = $1"

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

// DeleteTransactionByIDBulk handles DELETE /api/web/webpages/bulk/:id - DELETE
func DeleteTransactionByIDBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get ids array as a parameter as integer
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// Delete the webpage from the database
		for _, id := range ids {
			// query to delete the webpage
			query := "DELETE FROM transactions WHERE id = $1"

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

// UpdateTransactionStatusBulk handles PUT /api/web/webpages/status/bulk/:id - UPDATE
func UpdateTransactionStatusBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// Convert the string of ids to an array of ids
		ids := strings.Split(id, ",")

		// get the JSON data - only the status
		var Transaction models.WebpageModel
		if err := c.ShouldBindJSON(&Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the webpage status in the database
		for _, id := range ids {

			query := "UPDATE transactions SET status = $1 WHERE id = $2"

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

// UpdateTransactionStatus handles PUT /api/web/webpages/status/:id - UPDATE
func UpdateTransactionStatus(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get id parameter
		id := c.Param("id")

		// get the JSON data - only the status
		var Transaction models.WebpageModel
		if err := c.ShouldBindJSON(&Transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// query to update the webpage status
		query := "UPDATE transactions SET status = $1 WHERE id = $2"

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
