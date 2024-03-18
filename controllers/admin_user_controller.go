package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAdminUsers function
func GetAdminUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM admin_user")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into AdminUserModel structs
		var adminUsers []models.AdminUserModel

		for rows.Next() {
			var adminUser models.AdminUserModel
			if err := rows.Scan(&adminUser.ID, &adminUser.Name, &adminUser.Phone, &adminUser.Email, &adminUser.Password); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			adminUsers = append(adminUsers, adminUser)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all adminUsers as JSON
		c.JSON(http.StatusOK, adminUsers)

	}
}
