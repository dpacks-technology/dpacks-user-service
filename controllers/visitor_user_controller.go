package controllers

import (
	"database/sql"
	"dpacks-go-services-template/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetVisitorUsers function
func GetVisitorUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Query the database for all records
		rows, err := db.Query("SELECT * FROM visitor_user")

		if err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
			return
		}

		//close the rows when the surrounding function returns(handler function)
		defer rows.Close()

		// Iterate over the rows and scan them into VisitorUser structs
		var visitorUsers []models.VisitorUser

		for rows.Next() {
			var visitorUser models.VisitorUser
			if err := rows.Scan(&visitorUser.UserID, &visitorUser.Name, &visitorUser.Email, &visitorUser.PhoneNumber, &visitorUser.DateOfBirth, &visitorUser.Country, &visitorUser.FavoriteCategories, &visitorUser.UserDescription, &visitorUser.SignUpDate, &visitorUser.LastLogin, &visitorUser.ProfilePicture, &visitorUser.Gender, &visitorUser.Language, &visitorUser.Timezone); err != nil {
				fmt.Printf("%s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning rows from the database"})
				return
			}
			visitorUsers = append(visitorUsers, visitorUser)
		}

		//this runs only when loop didn't work
		if err := rows.Err(); err != nil {
			fmt.Printf("%s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over rows from the database"})
			return
		}

		// Return all visitorUsers as JSON
		c.JSON(http.StatusOK, visitorUsers)

	}
}
