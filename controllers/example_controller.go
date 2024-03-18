package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetExample handles GET /api/example - READ
func GetExample(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM example")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into ExampleModel structs
		var examples []models.ExampleModel

		for rows.Next() {
			var example models.ExampleModel
			if err := rows.Scan(&example.Column1, &example.Column2, &example.Column3); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			examples = append(examples, example)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all examples as JSON
		c.JSON(http.StatusOK, examples)

	}
}

// GetExampleByID handles GET /api/example/:id - READ
func GetExampleByID(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		// Create an empty ExampleModel struct
		var example models.ExampleModel

		// Query the database for the record with the given ID
		row := db.QueryRow("SELECT * FROM example WHERE new_column = $1", id)

		// Scan the row into the ExampleModel struct
		if err := row.Scan(&example.Column1, &example.Column2, &example.Column3); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, example)

	}
}

// AddExample handles POST /api/example - CREATE
func AddExample(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Create an empty ExampleModel struct
		var example models.ExampleModel

		// Bind the JSON to the ExampleModel struct
		if err := c.BindJSON(&example); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON"})
			return
		}

		// Insert the record into the database
		_, err := db.Exec("INSERT INTO example (column2, column3, new_column) VALUES ($1, $2, $3)", example.Column1, example.Column2, example.Column3)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting into the database"})
			return
		}

		// Return the example as JSON
		c.JSON(http.StatusOK, example)

	}
}

// UpdateExample handles PUT /api/example/:id - UPDATE
func UpdateExample(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		var example models.UpdateModel

		if err := c.BindJSON(&example); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Binding Data"})
			return
		}

		_, err := db.Exec("UPDATE example SET column2 = $1, column3 = $2 WHERE new_column=$3", example.Column1, example.Column2, id)
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in data Update"})
			return
		}
		// return statement
		c.JSON(http.StatusOK, gin.H{"Update": "Details Update"})
	}
}

// UpdateExampleBulk handles PUT /api/example/bulk - UPDATE
func UpdateExampleBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var example []models.ExampleModel

		//close the db connection

		if err := c.BindJSON(&example); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Data Binding Error"})
			return
		}

		//create tempory table and sned this data to tempory tbl
		_, err := db.Exec("CREATE TEMP TABLE temp_example (LIKE example INCLUDING ALL);")
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in creating tempory table"})
			return
		}

		//insert data into tempory table
		for _, v := range example {
			_, err := db.Exec("INSERT INTO temp_example (column2, column3, new_column) VALUES ($1, $2, $3)", v.Column1, v.Column2, v.Column3)
			if err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in inserting data into tempory table"})
				return
			}
		}

		//update data from tempory table to main table
		_, err = db.Exec("UPDATE example SET column2 = temp_example.column2, column3 = temp_example.column3 FROM temp_example WHERE example.new_column = temp_example.new_column")
		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in updating data"})
			return
		}

		// return statement
		c.JSON(http.StatusOK, gin.H{"success": "done"})
		db.Exec("DROP TABLE temp_example")
	}
}

// DeleteExample handles DELETE /api/example/:id - DELETE
func DeleteExample(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the ID from the URL
		id := c.Param("id")

		result, err := db.Exec("DELETE FROM example WHERE new_column = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete examples"})
			return
		}

		rowCount, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Deleted %d rows\n", rowCount)
		// return statement
		c.JSON(http.StatusOK, gin.H{"message": "Example deleted successfully"})

	}
}

// DeleteExampleBulk handles DELETE /api/example/bulk - DELETE
func DeleteExampleBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var request struct {
			IDs []int `json:"id"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		//Construct the DELETE query
		query := "DELETE FROM example WHERE new_column IN ("
		for i, id := range request.IDs {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("%d", id)
		}
		query += ")"

		result, err := db.Exec(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete examples"})
			return
		}

		rowCount, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Deleted %d rows\n", rowCount)

		// return statement
		c.JSON(http.StatusOK, gin.H{"message": "Example bulk deleted successfully"})
	}
}
